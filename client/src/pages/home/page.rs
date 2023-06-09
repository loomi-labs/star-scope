use sycamore::prelude::*;

use crate::config::keys;
use crate::utils::url::navigate_launch_app;

#[component]
pub fn LaunchButton<G: Html>(cx: Scope) -> View<G> {
    view!(cx,
        button(class="rounded-full flex items-center justify-center text-sm lg:text-2xl px-4 py-2 h-10 lg:h-12 w-36 lg:w-64 hover:text-black \
                transition-all bg-gradient-to-r from-primary-gradient-from to-primary-gradient-to hover:from-primary-gradient-to hover:to-primary-gradient-from",
                on:click=move |_| {
                    navigate_launch_app(cx)
                }) {
            span(class="") {
                "Launch App"
            }
            span(class="ml-2 icon-[pepicons-pop--arrow-right]") {}
        }
    )
}

#[component]
pub fn Intro<G: Html>(cx: Scope) -> View<G> {
    let title = "Bringing clarity to Cosmos";
    let description = "Star Scope is a notification service that helps you keep track of your activities across the Cosmos ecosystem.";

    view! {cx,
        div(class="hidden lg:flex justify-center min-h-screen") {
            div(class="flex flex-col justify-between 2xl:max-w-screen-2xl h-full w-full") {
                div(class="flex") {
                    div(class="flex flex-col justify-between w-1/3 min-h-screen pl-16 pt-16") {
                        img(class="h-auto w-64", src=keys::LOGO_WITH_TEXT_WHITE_IMG) {}
                        div(class="flex flex-col justify-start pb-16") {
                            h1(class="text-5xl font-bold") {
                                (title)
                            }
                            p(class="text-2xl my-8") {
                                (description)
                            }
                            LaunchButton {}
                        }
                        div {}
                    }
                    div(class="flex flex-col justify-center w-2/3 max-h-screen") {
                        img(class="h-full w-full object-contain", src=keys::SCOPE_IMG) {}
                    }
                }
            }
        }
        div(class="lg:hidden flex flex-col items-center min-h-screen p-8") {
            div(class="w-full") {
                div(class="flex justify-between items-center") {
                    img(class="h-auto w-32", src=keys::LOGO_WITH_TEXT_WHITE_IMG) {}
                    LaunchButton {}
                }
            }
            div(class="flex flex-col grow items-center justify-center text-center") {
                div(class="flex flex-col items-center py-16") {
                    h1(class="text-4xl font-bold") {
                        (title)
                    }
                    p(class="text-base my-8") {
                        (description)
                    }
                }
                div() {
                    img(class="h-full w-full object-contain", src=keys::SCOPE_IMG) {}
                }
            }
        }
    }
}

#[component]
pub fn Explanation<G: Html>(cx: Scope) -> View<G> {
    view! {cx,
         div(class="flex flex-col items-center justify-center min-h-screen p-8") {
             div(class="flex flex-col items-center justify-around min-h-screen 2xl:max-w-screen-2xl") {
                 h1(class="text-4xl md:text-5xl font-bold text-center py-16") {
                     "What is Star Scope?"
                 }
                 div(class="flex flex-col lg:flex-row justify-center items-center w-full pb-16") {
                     div(class="flex flex-col items-center justify-center w-full lg:w-1/3") {
                         div(class="flex flex-col items-center justify-center w-28 h-28 lg:w-56 lg:h-56 mt-8 mb-4 rounded-full bg-gradient-to-r from-purple-700 to-purple-900") {
                             div(class="w-16 h-16 lg:w-32 lg:h-32 icon-[mdi--email-fast-outline]") {}
                         }
                         h2(class="text-2xl font-bold my-4") {
                             "Instant Notifications"
                         }
                         p(class="text-center p-4") {
                             "Star Scope delivers real-time notifications about critical events occurring on Cosmos blockchains."
                         }
                     }
                     div(class="flex flex-col items-center justify-center w-full lg:w-1/3") {
                         div(class="flex flex-col items-center justify-center w-28 h-28 lg:w-56 lg:h-56 mt-8 mb-4 rounded-full bg-gradient-to-r from-purple-700 to-purple-900") {
                             div(class="w-16 h-16 lg:w-32 lg:h-32 icon-[octicon--bell-16]")
                         }
                         h2(class="text-2xl font-bold my-4") {
                             "Customizable Alerts"
                         }
                         p(class="text-center p-4") {
                             "Tailor your notification preferences to receive alerts specific to your interests."
                         }
                     }
                     div(class="flex flex-col items-center justify-center w-full lg:w-1/3") {
                         div(class="flex flex-col items-center justify-center w-28 h-28 lg:w-56 lg:h-56 mt-8 mb-4 rounded-full bg-gradient-to-r from-purple-700 to-purple-900") {
                             div(class="w-16 h-16 lg:w-32 lg:h-32 icon-[ps--world]")
                         }
                         h2(class="text-2xl font-bold my-4") {
                             "Extensive Coverage"
                         }
                         p(class="text-center p-4") {
                             "We provide comprehensive coverage of over 100 Cosmos blockchains."
                         }
                     }
                 }
                 div(class="mb-8") {
                     LaunchButton {}
                 }
            }
        }
    }
}

#[component]
pub fn Footer<G: Html>(cx: Scope) -> View<G> {
    view! {cx,
        div(class="flex flex-col items-center justify-center h-50 bg-d-bg-1") {
            div(class="flex w-full items-center justify-between max-w-screen-2xl p-8") {
                img(class="h-fit w-36 lg:w-64", src=keys::LOGO_WITH_TEXT_WHITE_IMG) {}
                div() {
                    a(class="p-4 hover:text-primary", href="https://t.me/rapha_decrypto", target="_blank") {
                        span(class="w-6 h-6 lg:w-14 lg:h-14 icon-[bxl--telegram]") {}
                    }
                    a(class="p-4 hover:text-primary", href="https://discord.com/users/228978159440232453", target="_blank") {
                        span(class="w-6 h-6 lg:w-14 lg:h-14 icon-[mingcute--discord-fill]") {}
                    }
                    a(class="p-4 hover:text-primary", href="https://twitter.com/StarScopeCosmos", target="_blank") {
                        span(class="w-6 h-6 lg:w-14 lg:h-14 icon-[mdi--twitter]") {}
                    }
                }
            }
        }
    }
}

#[component]
pub fn Home<G: Html>(cx: Scope) -> View<G> {
    view! {cx,
        div(class="bg-d-bg") {
            Intro {}
            Explanation {}
            Footer {}
        }
    }
}
