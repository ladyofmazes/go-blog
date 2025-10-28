package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/maxence-charriere/go-app/v10/pkg/ui"
)

const (
	headerHeight  = 72
	adsenseClient = "ca-pub-1013306768105236"
	adsenseSlot   = "9307554044"
)

type page struct {
	app.Compo

	Iclass   string
	Iindex   []app.UI
	Iicon    string
	Ititle   string
	Icontent []app.UI

	updateAvailable bool
	menuOpen        bool // Add this
}

func newPage() *page {
	return &page{}
}

func (p *page) Index(v ...app.UI) *page {
	p.Iindex = app.FilterUIElems(v...)
	return p
}

func (p *page) Icon(v string) *page {
	p.Iicon = v
	return p
}

func (p *page) Title(v string) *page {
	p.Ititle = v
	return p
}

func (p *page) Content(v ...app.UI) *page {
	p.Icontent = app.FilterUIElems(v...)
	return p
}

func (p *page) OnNav(ctx app.Context) {
	p.updateAvailable = ctx.AppUpdateAvailable()
	ctx.Defer(scrollTo)
}

func (p *page) OnAppUpdate(ctx app.Context) {
	p.updateAvailable = ctx.AppUpdateAvailable()
}

func (p *page) toggleMenu(ctx app.Context, e app.Event) {
	p.menuOpen = !p.menuOpen
}

func (p *page) Render() app.UI {
	shellClass := app.AppendClass("fill", "background")
	if p.menuOpen {
		shellClass = app.AppendClass(shellClass, "menu-open")
	}
	return ui.Shell().
		Class("fill").
		Class("background").
		Class(shellClass). // Add class when open
		Menu(
			newMenu().Class("fill"),
		).
		Index(
			app.If(len(p.Iindex) != 0, func() app.UI {
				return ui.Scroll().
					Class("fill").
					HeaderHeight(headerHeight).
					Content(
						app.Nav().
							Class("index").
							Body(
								app.Div().Class("separator"),
								app.Header().
									Class("h2").
									Text("Index"),
								app.Div().Class("separator"),
								app.Range(p.Iindex).Slice(func(i int) app.UI {
									return p.Iindex[i]
								}),
								newIndexLink().Title("Report an Issue"),
								app.Div().Class("separator"),
							),
					)
			}),
		).
		Content(
			ui.Scroll().
				Class("fill").
				Header(
					app.Nav().
						Class("fill").
						Body(
							ui.Stack().
								Class("fill").
								Left().
								Middle().
								Content(
									// Manual hamburger button
									app.Div().
										Class("hamburger-button").
										OnClick(p.toggleMenu).
										Body(
											app.Raw(`<svg viewBox="0 0 24 24" width="24" height="24">
                                                <path fill="currentColor" d="M3,6H21V8H3V6M3,11H21V13H3V11M3,16H21V18H3V16Z" />
                                            </svg>`),
										),

									app.If(p.updateAvailable, func() app.UI {
										return app.Div().
											Class("link-update").
											Body(
												ui.Link().
													Class("link").
													Class("heading").
													Class("fit").
													Class("unselectable").
													Icon(downloadSVG).
													Label("Update").
													OnClick(p.updateApp),
											)
									}),
								),
						),
				).
				HeaderHeight(headerHeight).
				Content(
					app.Main().Body(
						app.Article().Body(
							app.Header().
								ID("page-top").
								Class("page-title").
								Class("center").
								Body(
									ui.Stack().
										Center().
										Middle().
										Content(
											ui.Icon().
												Class("icon-left").
												Class("unselectable").
												Size(90).
												Src(p.Iicon),
											app.H1().Text(p.Ititle),
										),
								),
							app.Div().Class("separator"),
							app.Range(p.Icontent).Slice(func(i int) app.UI {
								return p.Icontent[i]
							}),

							app.Div().Class("separator"),
							app.Aside().Body(
								app.Header().
									ID("report-an-issue").
									Class("h2").
									Text(""),
								app.P().Body(
									app.Text("For more fun with me or to report an issue: "),
									app.A().
										Href("https://github.com/ladyofmazes/go-blog").
										Text("🚀 Find me on Github!"),
								),
							),
							app.Div().Class("separator"),
						),
					),
				),
		)
}

func (p *page) updateApp(ctx app.Context, e app.Event) {
	ctx.NewAction(updateApp)
}

func scrollTo(ctx app.Context) {
	id := ctx.Page().URL().Fragment
	if id == "" {
		id = "page-top"
	}
	ctx.ScrollTo(id)
}

func fragmentFocus(fragment string) string {
	if fragment == app.Window().URL().Fragment {
		return "focus"
	}
	return ""
}
