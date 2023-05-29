use sycamore::prelude::*;
use sycamore_router::navigate;

use crate::AppRoutes;
use crate::config::keys;
use crate::utils::url::safe_navigate;

#[component]
pub fn LaunchButton<G: Html>(cx: Scope) -> View<G> {
    view!(cx,
        button(class="rounded-full flex items-center justify-center text-sm lg:text-2xl px-4 py-2 h-10 lg:h-12 w-36 lg:w-64 hover:text-black \
                transition-all bg-gradient-to-r from-primary_gradient-from to-primary_gradient-to hover:from-primary_gradient-to hover:to-primary_gradient-from",
                on:click=move |_| safe_navigate(cx, AppRoutes::Login)) {
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
                        img(class="h-auto w-64", src=keys::LOGO_WITH_TEXT_IMG) {}
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
                    img(class="h-auto w-32", src=keys::LOGO_WITH_TEXT_IMG) {}
                    LaunchButton {}
                }
            }
            div(class="flex flex-col grow items-center justify-center text-center") {
                div(class="flex flex-col items-center py-16") {
                    h1(class="text-xl font-bold") {
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
        div(class="flex flex-col items-center justify-center min-h-screen") {
            div(class="flex flex-col items-center justify-around min-h-screen 2xl:max-w-screen-2xl") {
                h1(class="text-5xl font-bold text-center py-16") {
                    "What is Star Scope?"
                }
                div(class="flex flex-col lg:flex-row justify-center items-center w-full pb-16") {
                    div(class="flex flex-col items-center justify-center w-full lg:w-1/3") {
                        div(class="flex flex-col items-center justify-center w-28 h-28 lg:w-56 lg:h-56 mt-8 mb-4 rounded-full bg-gradient-to-r from-primary_gradient-from to-primary_gradient-to") {
                            div(class="w-16 h-16 lg:w-32 lg:h-32 icon-[mdi--email-fast-outline]") {}
                        }
                        h2(class="text-2xl font-bold my-4") {
                            "Instant Notifications"
                        }
                        p(class="text-center p-4") {
                            "Star Scope delivers real-time notifications about critical events occurring on Cosmos blockchains." br() "Stay informed and never miss important updates."
                        }
                    }
                    div(class="flex flex-col items-center justify-center w-full lg:w-1/3") {
                        div(class="flex flex-col items-center justify-center w-28 h-28 lg:w-56 lg:h-56 mt-8 mb-4 rounded-full bg-gradient-to-r from-primary_gradient-from to-primary_gradient-to") {
                            div(class="w-16 h-16 lg:w-32 lg:h-32 icon-[octicon--bell-16]")
                        }
                        h2(class="text-2xl font-bold my-4") {
                            "Customizable Alerts"
                        }
                        p(class="text-center p-4") {
                            "Tailor your notification preferences within Star Scope to receive alerts specific to your interests." br() "Choose the events and chains you want to monitor, enabling you to focus on what matters most to you."
                        }
                    }
                    div(class="flex flex-col items-center justify-center w-full lg:w-1/3") {
                        div(class="flex flex-col items-center justify-center w-28 h-28 lg:w-56 lg:h-56 mt-8 mb-4 rounded-full bg-gradient-to-r from-primary_gradient-from to-primary_gradient-to") {
                            div(class="w-16 h-16 lg:w-32 lg:h-32 icon-[ps--world]")
                        }
                        h2(class="text-2xl font-bold my-4") {
                            "Extensive Coverage"
                        }
                        p(class="text-center p-4") {
                            "We provide comprehensive coverage of over 100 Cosmos blockchains, ensuring you have access to vital information from the most prominent networks in the ecosystem." br() "Stay connected to the entire Cosmos community."
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
        div(class="flex flex-col items-center justify-center h-64 bg-landing_page-bg_footer") {
            div(class="flex w-full items-center justify-between max-w-screen-2xl") {
                img(class="h-fit w-36 lg:w-64", src=keys::LOGO_WITH_TEXT_IMG) {}
                div(class="") {
                    button(class="p-4 hover:text-primary", on:click=|_| navigate("https://t.me/rapha_decrypto")) {
                        span(class="w-6 h-6 lg:w-14 lg:h-14 icon-[bxl--telegram]") {}
                    }
                    button(class="p-4 hover:text-primary", on:click=|_| navigate("https://discord.com/users/228978159440232453")) {
                        span(class="w-6 h-6 lg:w-14 lg:h-14 icon-[mingcute--discord-fill]") {}
                    }
                    button(class="p-4 hover:text-primary", on:click=|_| navigate("https://twitter.com/Rapha90")) {
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
        div(class="bg-landing_page-bg") {
            Intro {}
            Explanation {}
            Footer {}
        }
    }
}
