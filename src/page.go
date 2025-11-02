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

	menu := app.Window().Get("document").Call("querySelector", ".menu")
	if menu.Truthy() {
		if p.menuOpen {
			style := menu.Get("style")
			style.Call("setProperty", "display", "block", "important")
			style.Call("setProperty", "position", "fixed", "important")
			style.Call("setProperty", "top", "0", "important")
			style.Call("setProperty", "left", "0", "important")
			style.Call("setProperty", "width", "100vw", "important") // Changed to 100vw
			style.Call("setProperty", "height", "100vh", "important")
			style.Call("setProperty", "z-index", "99999", "important")
			style.Call("setProperty", "background", "red", "important")
			style.Call("setProperty", "transform", "none", "important")
			style.Call("setProperty", "visibility", "visible", "important") // Add this
			style.Call("setProperty", "opacity", "1", "important")          // Add this

			app.Window().Get("console").Call("log", "Menu styles set with !important")

			// READ BACK the inline style to verify it stuck
			actualDisplay := style.Call("getPropertyValue", "display")
			actualBg := style.Call("getPropertyValue", "background")
			actualZIndex := style.Call("getPropertyValue", "z-index")

			app.Window().Get("console").Call("log", "Readback - display:", actualDisplay)
			app.Window().Get("console").Call("log", "Readback - background:", actualBg)
			app.Window().Get("console").Call("log", "Readback - z-index:", actualZIndex)

			// Also check the element's bounding box
			rect := menu.Call("getBoundingClientRect")
			app.Window().Get("console").Call("log", "Menu rect - top:", rect.Get("top"))
			app.Window().Get("console").Call("log", "Menu rect - left:", rect.Get("left"))
			app.Window().Get("console").Call("log", "Menu rect - width:", rect.Get("width"))
			app.Window().Get("console").Call("log", "Menu rect - height:", rect.Get("height"))
		}
	}
}

func (p *page) OnMount(ctx app.Context) {
	p.updateAvailable = ctx.AppUpdateAvailable()

	ctx.Async(func() {
		toggleFunc := app.FuncOf(func(this app.Value, args []app.Value) any {
			app.Window().Get("console").Call("log", "Creating custom menu overlay")

			if p.menuOpen {
				doc := app.Window().Get("document")
				body := doc.Get("body")

				// Create overlay
				overlay := doc.Call("createElement", "div")
				overlay.Set("id", "mobile-menu-overlay")

				// Copy the actual menu content from the .menu element
				originalMenu := doc.Call("querySelector", ".menu")
				if originalMenu.Truthy() {
					menuHTML := originalMenu.Get("innerHTML").String()
					overlay.Set("innerHTML", menuHTML)
				}

				style := overlay.Get("style")
				style.Set("position", "fixed")
				style.Set("top", "0")
				style.Set("left", "0")
				style.Set("width", "80%")
				style.Set("maxWidth", "300px")
				style.Set("height", "100vh")
				style.Set("zIndex", "99999")
				style.Set("background", "linear-gradient(#2e343a, rgba(0, 0, 0, 0.9))")
				style.Set("overflowY", "auto")
				style.Set("padding", "20px")

				body.Call("appendChild", overlay)
			} else {
				doc := app.Window().Get("document")
				overlay := doc.Call("getElementById", "mobile-menu-overlay")
				if overlay.Truthy() {
					overlay.Call("remove")
				}
			}

			return nil
		})

		app.Window().Set("toggleMenu", toggleFunc)
	})
}

func (p *page) Render() app.UI {
	shellClass := app.AppendClass("fill", "background")
	if p.menuOpen {
		shellClass = app.AppendClass(shellClass, "menu-open")
		shellClass = app.AppendClass(shellClass, "test-menu-is-open") // Extra debug class
	}
	return ui.Shell().
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
									// Manual hamburger button - REPLACE THIS ENTIRE SECTION
									app.Raw(`
                                        <div class="hamburger-button" onclick="window.toggleMenu()">
                                            <svg viewBox="0 0 24 24" width="24" height="24">
                                                <path fill="currentColor" d="M3,6H21V8H3V6M3,11H21V13H3V11M3,16H21V18H3V16Z" />
                                            </svg>
                                        </div>
                                    `),

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
