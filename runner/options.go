package runner

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/Nihility981/httpx_removebanner/common/customheader"
	"github.com/Nihility981/httpx_removebanner/common/customlist"
	customport "github.com/Nihility981/httpx_removebanner/common/customports"
	fileutilz "github.com/Nihility981/httpx_removebanner/common/fileutil"
	"github.com/Nihility981/httpx_removebanner/common/slice"
	"github.com/Nihility981/httpx_removebanner/common/stringz"
	"github.com/projectdiscovery/cdncheck"
	"github.com/projectdiscovery/fileutil"
	"github.com/projectdiscovery/goconfig"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/formatter"
	"github.com/projectdiscovery/gologger/levels"
)

const (
	// The maximum file length is 251 (255 - 4 bytes for ".ext" suffix)
	maxFileNameLength      = 251
	two                    = 2
	DefaultResumeFile      = "resume.cfg"
	DefaultOutputDirectory = "output"
)

var defaultProviders = strings.Join(cdncheck.GetDefaultProviders(), ", ")

type scanOptions struct {
	Methods                   []string
	StoreResponseDirectory    string
	RequestURI                string
	RequestBody               string
	VHost                     bool
	OutputTitle               bool
	OutputStatusCode          bool
	OutputLocation            bool
	OutputContentLength       bool
	StoreResponse             bool
	OutputServerHeader        bool
	OutputWebSocket           bool
	OutputWithNoColor         bool
	OutputMethod              bool
	ResponseInStdout          bool
	ChainInStdout             bool
	TLSProbe                  bool
	CSPProbe                  bool
	VHostInput                bool
	OutputContentType         bool
	Unsafe                    bool
	Pipeline                  bool
	HTTP2Probe                bool
	OutputIP                  bool
	OutputCName               bool
	OutputCDN                 bool
	OutputResponseTime        bool
	PreferHTTPS               bool
	NoFallback                bool
	NoFallbackScheme          bool
	TechDetect                bool
	StoreChain                bool
	MaxResponseBodySizeToSave int
	MaxResponseBodySizeToRead int
	OutputExtractRegex        string
	extractRegexps            map[string]*regexp.Regexp
	ExcludeCDN                bool
	HostMaxErrors             int
	ProbeAllIPS               bool
	Favicon                   bool
	LeaveDefaultPorts         bool
	OutputLinesCount          bool
	OutputWordsCount          bool
	Hashes                    string
}

func (s *scanOptions) Clone() *scanOptions {
	return &scanOptions{
		Methods:                   s.Methods,
		StoreResponseDirectory:    s.StoreResponseDirectory,
		RequestURI:                s.RequestURI,
		RequestBody:               s.RequestBody,
		VHost:                     s.VHost,
		OutputTitle:               s.OutputTitle,
		OutputStatusCode:          s.OutputStatusCode,
		OutputLocation:            s.OutputLocation,
		OutputContentLength:       s.OutputContentLength,
		StoreResponse:             s.StoreResponse,
		OutputServerHeader:        s.OutputServerHeader,
		OutputWebSocket:           s.OutputWebSocket,
		OutputWithNoColor:         s.OutputWithNoColor,
		OutputMethod:              s.OutputMethod,
		ResponseInStdout:          s.ResponseInStdout,
		ChainInStdout:             s.ChainInStdout,
		TLSProbe:                  s.TLSProbe,
		CSPProbe:                  s.CSPProbe,
		OutputContentType:         s.OutputContentType,
		Unsafe:                    s.Unsafe,
		Pipeline:                  s.Pipeline,
		HTTP2Probe:                s.HTTP2Probe,
		OutputIP:                  s.OutputIP,
		OutputCName:               s.OutputCName,
		OutputCDN:                 s.OutputCDN,
		OutputResponseTime:        s.OutputResponseTime,
		PreferHTTPS:               s.PreferHTTPS,
		NoFallback:                s.NoFallback,
		NoFallbackScheme:          s.NoFallbackScheme,
		TechDetect:                s.TechDetect,
		StoreChain:                s.StoreChain,
		OutputExtractRegex:        s.OutputExtractRegex,
		MaxResponseBodySizeToSave: s.MaxResponseBodySizeToSave,
		MaxResponseBodySizeToRead: s.MaxResponseBodySizeToRead,
		HostMaxErrors:             s.HostMaxErrors,
		Favicon:                   s.Favicon,
		extractRegexps:            s.extractRegexps,
		LeaveDefaultPorts:         s.LeaveDefaultPorts,
		OutputLinesCount:          s.OutputLinesCount,
		OutputWordsCount:          s.OutputWordsCount,
		Hashes:                    s.Hashes,
	}
}

// Options contains configuration options for httpx.
type Options struct {
	CustomHeaders             customheader.CustomHeaders
	CustomPorts               customport.CustomPorts
	matchStatusCode           []int
	matchContentLength        []int
	filterStatusCode          []int
	filterContentLength       []int
	Output                    string
	StoreResponseDir          string
	HTTPProxy                 string
	SocksProxy                string
	InputFile                 string
	Methods                   string
	RequestURI                string
	RequestURIs               string
	requestURIs               []string
	OutputMatchStatusCode     string
	OutputMatchContentLength  string
	OutputFilterStatusCode    string
	OutputFilterContentLength string
	InputRawRequest           string
	rawRequest                string
	RequestBody               string
	OutputFilterString        string
	OutputMatchString         string
	OutputFilterRegex         string
	OutputMatchRegex          string
	Retries                   int
	Threads                   int
	Timeout                   int
	filterRegex               *regexp.Regexp
	matchRegex                *regexp.Regexp
	VHost                     bool
	VHostInput                bool
	Smuggling                 bool
	ExtractTitle              bool
	StatusCode                bool
	Location                  bool
	ContentLength             bool
	FollowRedirects           bool
	StoreResponse             bool
	JSONOutput                bool
	CSVOutput                 bool
	Silent                    bool
	Version                   bool
	Verbose                   bool
	NoColor                   bool
	OutputServerHeader        bool
	OutputWebSocket           bool
	responseInStdout          bool
	chainInStdout             bool
	FollowHostRedirects       bool
	MaxRedirects              int
	OutputMethod              bool
	TLSProbe                  bool
	CSPProbe                  bool
	OutputContentType         bool
	OutputIP                  bool
	OutputCName               bool
	Unsafe                    bool
	Debug                     bool
	DebugRequests             bool
	DebugResponse             bool
	Pipeline                  bool
	HTTP2Probe                bool
	OutputCDN                 bool
	OutputResponseTime        bool
	NoFallback                bool
	NoFallbackScheme          bool
	TechDetect                bool
	TLSGrab                   bool
	protocol                  string
	ShowStatistics            bool
	StatsInterval             int
	RandomAgent               bool
	StoreChain                bool
	Deny                      customlist.CustomList
	Allow                     customlist.CustomList
	MaxResponseBodySizeToSave int
	MaxResponseBodySizeToRead int
	OutputExtractRegexs       goflags.StringSlice
	OutputExtractPresets      goflags.StringSlice
	RateLimit                 int
	RateLimitMinute           int
	Probe                     bool
	Resume                    bool
	resumeCfg                 *ResumeCfg
	ExcludeCDN                bool
	HostMaxErrors             int
	Stream                    bool
	SkipDedupe                bool
	ProbeAllIPS               bool
	Resolvers                 goflags.StringSlice
	Favicon                   bool
	OutputFilterFavicon       goflags.StringSlice
	OutputMatchFavicon        goflags.StringSlice
	LeaveDefaultPorts         bool
	OutputLinesCount          bool
	OutputMatchLinesCount     string
	matchLinesCount           []int
	OutputFilterLinesCount    string
	filterLinesCount          []int
	OutputWordsCount          bool
	OutputMatchWordsCount     string
	matchWordsCount           []int
	OutputFilterWordsCount    string
	filterWordsCount          []int
	Hashes                    string
	Jarm                      bool
	Asn                       bool
	OutputMatchCdn            goflags.StringSlice
	OutputFilterCdn           goflags.StringSlice
	SniName                   string
	OutputMatchResponseTime   string
	OutputFilterResponseTime  string
	HealthCheck               bool
}

// ParseOptions parses the command line options for application
func ParseOptions() *Options {
	options := &Options{}

	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`httpx is a fast and multi-purpose HTTP toolkit allow to run multiple probers using retryablehttp library.`)

	flagSet.CreateGroup("input", "Input",
		flagSet.StringVarP(&options.InputFile, "list", "l", "", "input file containing list of hosts to process"),
		flagSet.StringVarP(&options.InputRawRequest, "request", "rr", "", "file containing raw request"),
	)

	flagSet.CreateGroup("Probes", "Probes",
		flagSet.BoolVarP(&options.StatusCode, "status-code", "sc", false, "display response status-code"),
		flagSet.BoolVarP(&options.ContentLength, "content-length", "cl", false, "display response content-length"),
		flagSet.BoolVarP(&options.OutputContentType, "content-type", "ct", false, "display response content-type"),
		flagSet.BoolVar(&options.Location, "location", false, "display response redirect location"),
		flagSet.BoolVar(&options.Favicon, "favicon", false, "display mmh3 hash for '/favicon.ico' file"),
		flagSet.StringVar(&options.Hashes, "hash", "", "display response body hash (supported: md5,mmh3,simhash,sha1,sha256,sha512)"),
		flagSet.BoolVar(&options.Jarm, "jarm", false, "display jarm fingerprint hash"),
		flagSet.BoolVarP(&options.OutputResponseTime, "response-time", "rt", false, "display response time"),
		flagSet.BoolVarP(&options.OutputLinesCount, "line-count", "lc", false, "display response body line count"),
		flagSet.BoolVarP(&options.OutputWordsCount, "word-count", "wc", false, "display response body word count"),
		flagSet.BoolVar(&options.ExtractTitle, "title", false, "display page title"),
		flagSet.BoolVarP(&options.OutputServerHeader, "web-server", "server", false, "display server name"),
		flagSet.BoolVarP(&options.TechDetect, "tech-detect", "td", false, "display technology in use based on wappalyzer dataset"),
		flagSet.BoolVar(&options.OutputMethod, "method", false, "display http request method"),
		flagSet.BoolVar(&options.OutputWebSocket, "websocket", false, "display server using websocket"),
		flagSet.BoolVar(&options.OutputIP, "ip", false, "display host ip"),
		flagSet.BoolVar(&options.OutputCName, "cname", false, "display host cname"),
		flagSet.BoolVar(&options.Asn, "asn", false, "display host asn information"),
		flagSet.BoolVar(&options.OutputCDN, "cdn", false, "display cdn in use"),
		flagSet.BoolVar(&options.Probe, "probe", false, "display probe status"),
	)

	flagSet.CreateGroup("matchers", "Matchers",
		flagSet.StringVarP(&options.OutputMatchStatusCode, "match-code", "mc", "", "match response with specified status code (-mc 200,302)"),
		flagSet.StringVarP(&options.OutputMatchContentLength, "match-length", "ml", "", "match response with specified content length (-ml 100,102)"),
		flagSet.StringVarP(&options.OutputMatchLinesCount, "match-line-count", "mlc", "", "match response body with specified line count (-mlc 423,532)"),
		flagSet.StringVarP(&options.OutputMatchWordsCount, "match-word-count", "mwc", "", "match response body with specified word count (-mwc 43,55)"),
		flagSet.StringSliceVarP(&options.OutputMatchFavicon, "match-favicon", "mfc", nil, "match response with specified favicon hash (-mfc 1494302000)", goflags.NormalizedStringSliceOptions),
		flagSet.StringVarP(&options.OutputMatchString, "match-string", "ms", "", "match response with specified string (-ms admin)"),
		flagSet.StringVarP(&options.OutputMatchRegex, "match-regex", "mr", "", "match response with specified regex (-mr admin)"),
		flagSet.StringSliceVarP(&options.OutputMatchCdn, "match-cdn", "mcdn", nil, fmt.Sprintf("match host with specified cdn provider (%s)", defaultProviders), goflags.NormalizedStringSliceOptions),
		flagSet.StringVarP(&options.OutputMatchResponseTime, "match-response-time", "mrt", "", "match response with specified response time in seconds (-mrt '< 1')"),
	)

	flagSet.CreateGroup("extractor", "Extractor",
		flagSet.StringSliceVarP(&options.OutputExtractRegexs, "extract-regex", "er", nil, "Display response content with matched regex", goflags.StringSliceOptions),
		flagSet.StringSliceVarP(&options.OutputExtractPresets, "extract-preset", "ep", nil, "Display response content with matched preset regex", goflags.StringSliceOptions),
	)

	flagSet.CreateGroup("filters", "Filters",
		flagSet.StringVarP(&options.OutputFilterStatusCode, "filter-code", "fc", "", "filter response with specified status code (-fc 403,401)"),
		flagSet.StringVarP(&options.OutputFilterContentLength, "filter-length", "fl", "", "filter response with specified content length (-fl 23,33)"),
		flagSet.StringVarP(&options.OutputFilterLinesCount, "filter-line-count", "flc", "", "filter response body with specified line count (-flc 423,532)"),
		flagSet.StringVarP(&options.OutputFilterWordsCount, "filter-word-count", "fwc", "", "filter response body with specified word count (-fwc 423,532)"),
		flagSet.StringSliceVarP(&options.OutputFilterFavicon, "filter-favicon", "ffc", nil, "filter response with specified favicon hash (-mfc 1494302000)", goflags.NormalizedStringSliceOptions),
		flagSet.StringVarP(&options.OutputFilterString, "filter-string", "fs", "", "filter response with specified string (-fs admin)"),
		flagSet.StringVarP(&options.OutputFilterRegex, "filter-regex", "fe", "", "filter response with specified regex (-fe admin)"),
		flagSet.StringSliceVarP(&options.OutputFilterCdn, "filter-cdn", "fcdn", nil, fmt.Sprintf("filter host with specified cdn provider (%s)", defaultProviders), goflags.NormalizedStringSliceOptions),
		flagSet.StringVarP(&options.OutputFilterResponseTime, "filter-response-time", "frt", "", "filter response with specified response time in seconds (-frt '> 1')"),
	)

	flagSet.CreateGroup("rate-limit", "Rate-Limit",
		flagSet.IntVarP(&options.Threads, "threads", "t", 50, "number of threads to use"),
		flagSet.IntVarP(&options.RateLimit, "rate-limit", "rl", 150, "maximum requests to send per second"),
		flagSet.IntVarP(&options.RateLimitMinute, "rate-limit-minute", "rlm", 0, "maximum number of requests to send per minute"),
	)

	flagSet.CreateGroup("Misc", "Miscellaneous",
		flagSet.BoolVarP(&options.ProbeAllIPS, "probe-all-ips", "pa", false, "probe all the ips associated with same host"),
		flagSet.VarP(&options.CustomPorts, "ports", "p", "ports to probe (nmap syntax: eg http:1,2-10,11,https:80)"),
		flagSet.StringVar(&options.RequestURIs, "path", "", "path or list of paths to probe (comma-separated, file)"),
		flagSet.BoolVar(&options.TLSProbe, "tls-probe", false, "send http probes on the extracted TLS domains (dns_name)"),
		flagSet.BoolVar(&options.CSPProbe, "csp-probe", false, "send http probes on the extracted CSP domains"),
		flagSet.BoolVar(&options.TLSGrab, "tls-grab", false, "perform TLS(SSL) data grabbing"),
		flagSet.BoolVar(&options.Pipeline, "pipeline", false, "probe and display server supporting HTTP1.1 pipeline"),
		flagSet.BoolVar(&options.HTTP2Probe, "http2", false, "probe and display server supporting HTTP2"),
		flagSet.BoolVar(&options.VHost, "vhost", false, "probe and display server supporting VHOST"),
	)

	flagSet.CreateGroup("output", "Output",
		flagSet.StringVarP(&options.Output, "output", "o", "", "file to write output results"),
		flagSet.BoolVarP(&options.StoreResponse, "store-response", "sr", false, "store http response to output directory"),
		flagSet.StringVarP(&options.StoreResponseDir, "store-response-dir", "srd", "", "store http response to custom directory"),
		flagSet.BoolVar(&options.CSVOutput, "csv", false, "store output in csv format"),
		flagSet.BoolVar(&options.JSONOutput, "json", false, "store output in JSONL(ines) format"),
		flagSet.BoolVarP(&options.responseInStdout, "include-response", "irr", false, "include http request/response in JSON output (-json only)"),
		flagSet.BoolVar(&options.chainInStdout, "include-chain", false, "include redirect http chain in JSON output (-json only)"),
		flagSet.BoolVar(&options.StoreChain, "store-chain", false, "include http redirect chain in responses (-sr only)"),
	)

	flagSet.CreateGroup("configs", "Configurations",
		flagSet.StringSliceVarP(&options.Resolvers, "resolvers", "r", nil, "list of custom resolver (file or comma separated)", goflags.NormalizedStringSliceOptions),
		flagSet.Var(&options.Allow, "allow", "allowed list of IP/CIDR's to process (file or comma separated)"),
		flagSet.Var(&options.Deny, "deny", "denied list of IP/CIDR's to process (file or comma separated)"),
		flagSet.StringVarP(&options.SniName, "sni-name", "sni", "", "Custom TLS SNI name"),
		flagSet.BoolVar(&options.RandomAgent, "random-agent", true, "Enable Random User-Agent to use"),
		flagSet.VarP(&options.CustomHeaders, "header", "H", "custom http headers to send with request"),
		flagSet.StringVarP(&options.HTTPProxy, "proxy", "http-proxy", "", "http proxy to use (eg http://127.0.0.1:8080)"),
		flagSet.BoolVar(&options.Unsafe, "unsafe", false, "send raw requests skipping golang normalization"),
		flagSet.BoolVar(&options.Resume, "resume", false, "resume scan using resume.cfg"),
		flagSet.BoolVarP(&options.FollowRedirects, "follow-redirects", "fr", false, "follow http redirects"),
		flagSet.IntVarP(&options.MaxRedirects, "max-redirects", "maxr", 10, "max number of redirects to follow per host"),
		flagSet.BoolVarP(&options.FollowHostRedirects, "follow-host-redirects", "fhr", false, "follow redirects on the same host"),
		flagSet.BoolVar(&options.VHostInput, "vhost-input", false, "get a list of vhosts as input"),
		flagSet.StringVar(&options.Methods, "x", "", "request methods to probe, use 'all' to probe all HTTP methods"),
		flagSet.StringVar(&options.RequestBody, "body", "", "post body to include in http request"),
		flagSet.BoolVarP(&options.Stream, "stream", "s", false, "stream mode - start elaborating input targets without sorting"),
		flagSet.BoolVarP(&options.SkipDedupe, "skip-dedupe", "sd", false, "disable dedupe input items (only used with stream mode)"),
		flagSet.BoolVarP(&options.LeaveDefaultPorts, "leave-default-ports", "ldp", false, "leave default http/https ports in host header (eg. http://host:80 - https//host:443"),
	)

	flagSet.CreateGroup("debug", "Debug",
		flagSet.BoolVarP(&options.HealthCheck, "hc", "health-check", false, "run diagnostic check up"),
		flagSet.BoolVar(&options.Debug, "debug", false, "display request/response content in cli"),
		flagSet.BoolVar(&options.DebugRequests, "debug-req", false, "display request content in cli"),
		flagSet.BoolVar(&options.DebugResponse, "debug-resp", false, "display response content in cli"),
		flagSet.BoolVar(&options.Version, "version", false, "display httpx version"),
		flagSet.BoolVar(&options.ShowStatistics, "stats", false, "display scan statistic"),
		flagSet.BoolVar(&options.Silent, "silent", false, "silent mode"),
		flagSet.BoolVarP(&options.Verbose, "verbose", "v", false, "verbose mode"),
		flagSet.IntVarP(&options.StatsInterval, "stats-interval", "si", 0, "number of seconds to wait between showing a statistics update (default: 5)"),
		flagSet.BoolVarP(&options.NoColor, "no-color", "nc", false, "disable colors in cli output"),
	)

	flagSet.CreateGroup("Optimizations", "Optimizations",
		flagSet.BoolVarP(&options.NoFallback, "no-fallback", "nf", false, "display both probed protocol (HTTPS and HTTP)"),
		flagSet.BoolVarP(&options.NoFallbackScheme, "no-fallback-scheme", "nfs", false, "probe with protocol scheme specified in input "),
		flagSet.IntVarP(&options.HostMaxErrors, "max-host-error", "maxhr", 30, "max error count per host before skipping remaining path/s"),
		flagSet.BoolVarP(&options.ExcludeCDN, "exclude-cdn", "ec", false, "skip full port scans for CDNs (only checks for 80,443)"),
		flagSet.IntVar(&options.Retries, "retries", 0, "number of retries"),
		flagSet.IntVar(&options.Timeout, "timeout", 5, "timeout in seconds"),
		flagSet.IntVarP(&options.MaxResponseBodySizeToSave, "response-size-to-save", "rsts", math.MaxInt32, "max response size to save in bytes"),
		flagSet.IntVarP(&options.MaxResponseBodySizeToRead, "response-size-to-read", "rstr", math.MaxInt32, "max response size to read in bytes"),
	)

	_ = flagSet.Parse()

	if options.HealthCheck {
		gologger.Print().Msgf("%s\n", DoHealthCheck(options))
		os.Exit(0)
	}

	if options.StatsInterval != 0 {
		options.ShowStatistics = true
	}

	// Read the inputs and configure the logging
	options.configureOutput()

	err := options.configureResume()
	if err != nil {
		gologger.Fatal().Msgf("%s\n", err)
	}

	showBanner()

	if options.Version {
		gologger.Info().Msgf("Current Version: %s\n", Version)
		os.Exit(0)
	}

	if err := options.ValidateOptions(); err != nil {
		gologger.Fatal().Msgf("%s\n", err)
	}

	return options
}

func (options *Options) ValidateOptions() error {
	if options.InputFile != "" && !fileutilz.FileNameIsGlob(options.InputFile) && !fileutil.FileExists(options.InputFile) {
		return fmt.Errorf("File %s does not exist.", options.InputFile)
	}

	if options.InputRawRequest != "" && !fileutil.FileExists(options.InputRawRequest) {
		return fmt.Errorf("File %s does not exist.", options.InputRawRequest)
	}

	multiOutput := options.CSVOutput && options.JSONOutput
	if multiOutput {
		return fmt.Errorf("Results can only be displayed in one format: 'JSON' or 'CSV'")
	}

	var err error
	if options.matchStatusCode, err = stringz.StringToSliceInt(options.OutputMatchStatusCode); err != nil {
		return errors.Wrap(err, "Invalid value for match status code option")
	}
	if options.matchContentLength, err = stringz.StringToSliceInt(options.OutputMatchContentLength); err != nil {
		return errors.Wrap(err, "Invalid value for match content length option")
	}
	if options.filterStatusCode, err = stringz.StringToSliceInt(options.OutputFilterStatusCode); err != nil {
		return errors.Wrap(err, "Invalid value for filter status code option")
	}
	if options.filterContentLength, err = stringz.StringToSliceInt(options.OutputFilterContentLength); err != nil {
		return errors.Wrap(err, "Invalid value for filter content length option")
	}
	if options.OutputFilterRegex != "" {
		if options.filterRegex, err = regexp.Compile(options.OutputFilterRegex); err != nil {
			return errors.Wrap(err, "Invalid value for regex filter option")
		}
	}
	if options.OutputMatchRegex != "" {
		if options.matchRegex, err = regexp.Compile(options.OutputMatchRegex); err != nil {
			return errors.Wrap(err, "Invalid value for match regex option")
		}
	}
	if options.matchLinesCount, err = stringz.StringToSliceInt(options.OutputMatchLinesCount); err != nil {
		return errors.Wrap(err, "Invalid value for match lines count option")
	}
	if options.matchWordsCount, err = stringz.StringToSliceInt(options.OutputMatchWordsCount); err != nil {
		return errors.Wrap(err, "Invalid value for match words count option")
	}
	if options.filterLinesCount, err = stringz.StringToSliceInt(options.OutputFilterLinesCount); err != nil {
		return errors.Wrap(err, "Invalid value for filter lines count option")
	}
	if options.filterWordsCount, err = stringz.StringToSliceInt(options.OutputFilterWordsCount); err != nil {
		return errors.Wrap(err, "Invalid value for filter words count option")
	}

	var resolvers []string
	for _, resolver := range options.Resolvers {
		if fileutil.FileExists(resolver) {
			chFile, err := fileutil.ReadFile(resolver)
			if err != nil {
				return errors.Wrapf(err, "Couldn't process resolver file \"%s\"", resolver)
			}
			for line := range chFile {
				resolvers = append(resolvers, line)
			}
		} else {
			resolvers = append(resolvers, resolver)
		}
	}

	options.Resolvers = resolvers
	if len(options.Resolvers) > 0 {
		gologger.Debug().Msgf("Using resolvers: %s\n", strings.Join(options.Resolvers, ","))
	}

	if options.StoreResponse && options.StoreResponseDir == "" {
		gologger.Debug().Msgf("Store response directory not specified, using \"%s\"\n", DefaultOutputDirectory)
		options.StoreResponseDir = DefaultOutputDirectory
	}
	if options.StoreResponseDir != "" && !options.StoreResponse {
		gologger.Debug().Msgf("Store response directory specified, enabling \"sr\" flag automatically\n")
		options.StoreResponse = true
	}

	if options.Favicon {
		gologger.Debug().Msgf("Setting single path to \"favicon.ico\" and ignoring multiple paths settings\n")
		options.RequestURIs = "/favicon.ico"
	}

	if options.Hashes != "" {
		for _, hashType := range strings.Split(options.Hashes, ",") {
			if !slice.StringSliceContains([]string{"md5", "sha1", "sha256", "sha512", "mmh3", "simhash"}, strings.ToLower(hashType)) {
				gologger.Error().Msgf("Unsupported hash type: %s\n", hashType)
			}
		}
	}
	if len(options.OutputMatchCdn) > 0 || len(options.OutputFilterCdn) > 0 {
		options.OutputCDN = true
	}

	return nil
}

// configureOutput configures the output on the screen
func (options *Options) configureOutput() {
	// If the user desires verbose output, show verbose output
	if options.Verbose {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelVerbose)
	}
	if options.Debug {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
	}
	if options.NoColor {
		gologger.DefaultLogger.SetFormatter(formatter.NewCLI(true))
	}
	if options.Silent {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelSilent)
	}
	if len(options.OutputMatchResponseTime) > 0 || len(options.OutputFilterResponseTime) > 0 {
		options.OutputResponseTime = true
	}
}

func (options *Options) configureResume() error {
	options.resumeCfg = &ResumeCfg{}
	if options.Resume && fileutil.FileExists(DefaultResumeFile) {
		return goconfig.Load(&options.resumeCfg, DefaultResumeFile)

	}
	return nil
}

// ShouldLoadResume resume file
func (options *Options) ShouldLoadResume() bool {
	return options.Resume && fileutil.FileExists(DefaultResumeFile)
}

// ShouldSaveResume file
func (options *Options) ShouldSaveResume() bool {
	return true
}
