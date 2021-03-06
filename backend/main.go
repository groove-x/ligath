package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

type TestableHandler func(w http.ResponseWriter, r *http.Request, isMock bool)
type Middleware func(next http.Handler) http.Handler

type jsons []string

func (j *jsons) String() string {
	return "JSONs"
}

func (j *jsons) Set(value string) error {
	*j = append(*j, value)
	return nil
}

func (j *jsons) Array() []string {
	return *j
}

var helpLines = []string{
	fmt.Sprintf("%s <help|serve|export> [OPTION]...", os.Args[0]),
	"Subcommands:",
	"    help     Show this help",
	"    serve    Serve Ligath frontend & backend",
	"    export   Export verified package information",
	"Options:",
	"    -h       Show this help",
}

func replaceUsage(fs *flag.FlagSet) {
	usage := func() {
		fmt.Println(strings.Join(helpLines, "\n"))
	}

	if fs == nil {
		flag.Usage = usage
	} else {
		fs.Usage = usage
	}
}

func main() {
	jsons := &jsons{}

	replaceUsage(nil)
	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	switch flag.Arg(0) {
	default:
		fmt.Println("Unknown subcommand: " + flag.Arg(0))
		fallthrough
	case "", "help":
		flag.Usage()
		return

	case "serve":
		helpLines = []string{
			fmt.Sprintf("%s serve [OPTION]...", os.Args[0]),
			"Options:",
			"    -j <JSONPATH>  JSON to migrate",
			"Examples:",
			"    ligath serve",
			"    ligath serve -j debian.json",
			"    ligath serve -j debian_apt.json -j ubuntu_apt.json",
		}

		serve := flag.NewFlagSet("serve", flag.ExitOnError)
		replaceUsage(serve)

		serve.Var(jsons, "j", "JSON to migrate")
		serve.Parse(os.Args[2:])

		a := NewAPI("ligath.db", jsons.Array())
		defer a.Close()

		err := a.Serve()
		if err != nil {
			panic(err)
		}

	case "export":
		helpLines = []string{
			fmt.Sprintf("%s export <FORMAT> <OUTFILE>", os.Args[0]),
			"Arguments:",
			"    FORMAT   Output format (available choice: `human`)",
			"    OUTFILE  Output file name (`-` for stdout)",
		}

		export := flag.NewFlagSet("export", flag.ExitOnError)
		replaceUsage(export)
		export.Parse(os.Args[2:])

		if export.NArg() < 2 {
			export.Usage()
			return
		}

		format, outfile := export.Arg(0), export.Arg(1)
		a := NewAPI("ligath.db", []string{})
		defer a.Close()

		a.export(format, outfile)
	}
}

type API struct {
	router  *chi.Mux
	storage Storage
}

func NewAPI(dbPath string, jsons []string) *API {
	a := &API{}
	r := chi.NewRouter()

	val, ok := os.LookupEnv("MOCK")
	isMock := ok && val != ""
	r.Use(setMock(isMock))
	r.Use(middleware.Logger)

	cross := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	})
	r.Use(cross.Handler)

	r.Route("/api", func(r chi.Router) {
		r.Get("/packages", toHandlerFunc(a.getPackages))
		r.Get("/packages/{package}@{version}", toHandlerFunc(a.getPackage))
		r.Put("/packages/{package}@{version}", toHandlerFunc(a.putPackage))
		r.Get("/licenses", toHandlerFunc(a.getLicenses))
	})

	var s Storage
	if isMock {
		s = NewMockStorage()
	} else {
		var err error
		s, err = NewBoltStorage(dbPath, jsons)
		if err != nil {
			panic(err)
		}
	}

	wd, _ := os.Getwd()
	if !strings.HasSuffix(wd, "ligath") {
		fmt.Println("Please run it in repository root as working directory.")
		fmt.Printf("Current: %s\n", wd)
		os.Exit(1)
	}

	r.Get("/css/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(wd, "frontend/dist/", r.URL.Path[1:]))
	}))
	r.Get("/js/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(wd, "frontend/dist/", r.URL.Path[1:]))
	}))
	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(wd, "frontend/dist/index.html"))
	}))

	a.router = r
	a.storage = s
	return a
}

func (a *API) Close() {
	a.storage.Close()
}

func (a *API) Serve() error {
	err := http.ListenAndServe("0.0.0.0:3939", a.router)
	if err != nil {
		return fmt.Errorf("error while serving: %v", err)
	}
	return nil
}

func (a *API) getPackage(w http.ResponseWriter, r *http.Request, isMock bool) {
	var err error
	var pkg *Package

	paramPkg := chi.URLParam(r, "package")
	paramVer := chi.URLParam(r, "version")
	paramKind := r.URL.Query().Get("kind")

	if paramPkg != "" && paramVer != "" {
		if paramKind == "" {
			pkg, err = a.storage.GetPackage(paramPkg, paramVer)
		} else if paramKind == "parsed" {
			pkg, err = a.storage.GetParsedPackage(paramPkg, paramVer)
		} else if paramKind == "notparsed" {
			pkg, err = a.storage.GetNotParsedPackage(paramPkg, paramVer)
		} else if paramKind == "verified" {
			pkg, err = a.storage.GetVerifiedPackage(paramPkg, paramVer)
		}

		if err != nil {
			// TODO: handle error
			render.Status(r, 500)
			render.PlainText(w, r, "")
		} else if pkg == nil {
			render.Status(r, 404)
			render.PlainText(w, r, "")
		} else {
			render.JSON(w, r, pkg)
		}
	} else {
		render.Status(r, 400)
		render.PlainText(w, r, "")
	}
}

func (a *API) putPackage(w http.ResponseWriter, r *http.Request, isMock bool) {
	pkg := chi.URLParam(r, "package")
	ver := chi.URLParam(r, "version")
	if pkg == "" || ver == "" {
		render.Status(r, 400)
		return
	}

	j := Package{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&j)
	if err != nil {
		log.Printf("failed to decode JSON: %v", err)
		render.Status(r, 400)
	}

	err = a.storage.PutPackage(j)
	if err != nil {
		log.Printf("failed to put package: %v", err)
		render.Status(r, 500)
	}

	render.Status(r, 200)
}

func (a *API) getPackages(w http.ResponseWriter, r *http.Request, isMock bool) {
	var pkgs []PackageListItem

	if license := r.URL.Query().Get("license"); license != "" {
		pkgs = a.storage.GetPackagesWithLicense(license)
	} else if kind := r.URL.Query().Get("kind"); kind != "" {
		switch kind {
		case "parsed":
			pkgs = a.storage.GetParsedPackages()
		case "notparsed":
			pkgs = a.storage.GetNotParsedPackages()
		case "verified":
			pkgs = a.storage.GetVerifiedPackages()
		case "emptycopyright":
			pkgs = a.storage.GetEmptyCopyrightPackages()
		default:
			render.Status(r, 400)
		}
	}

	if len(pkgs) > 0 {
		render.JSON(w, r, pkgs)
	} else {
		render.Status(r, 404)
	}
}

func (a *API) getLicenses(w http.ResponseWriter, r *http.Request, isMock bool) {
	render.JSON(w, r, a.storage.GetLicenses())
}

func (a *API) export(format, filename string) {
	var err error
	var out *os.File

	if filename == "-" {
		out = os.Stdout
	} else {
		out, err = os.Create(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create output file %s: %s\n", filename, err)
			return
		}
		defer out.Close()
	}

	var text string

	switch format {
	case "human":
		text = a.exportHumanReadable()
	default:
		fmt.Fprintln(os.Stderr, "Unknown format: "+format)
	}

	out.WriteString(text)
}

func (a *API) exportHumanReadable() string {
	var buf []string

	packages := a.storage.GetVerifiedPackages()
	for _, pkg := range packages {
		info, err := a.storage.GetVerifiedPackage(pkg.Name, pkg.Version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get package %s@%s: %s\n", pkg.Name, pkg.Version, err)
			continue
		}

		delim := strings.Repeat("=", len(info.Name)+len(info.Version)+1)
		s := fmt.Sprintf(
			"%s\n%s %s\n%s\n\n%s",
			delim,
			pkg.Name,
			pkg.Version,
			delim,
			info.RawCopyright,
		)
		buf = append(buf, s)
	}

	return strings.Join(buf, "\n\n")
}

func setMock(mock bool) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "mock", mock))
			next.ServeHTTP(w, r)
		})
	}
}

func toHandlerFunc(handler TestableHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		isMock, ok := r.Context().Value("mock").(bool)
		if !ok {
			panic("illegal mock flag is passed")
		}

		handler(w, r, isMock)
	}
}
