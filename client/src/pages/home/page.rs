use sycamore::prelude::*;
use crate::AppRoutes;
use crate::config::keys;
use crate::utils::url::safe_navigate;

#[component]
pub fn LaunchButton<G: Html>(cx: Scope) -> View<G> {
    view!(cx,
        button(class="rounded-full flex items-center justify-center px-4 py-2 h-10 lg:h-12 w-48 lg:w-64 text-black hover:text-white bg-gradient-to-r from-primary_gradient-from to-primary_gradient-to",
                on:click=move |_| safe_navigate(cx, AppRoutes::Login)) {
            span(class="text-xl") {
                "Launch App"
            }
            span(class="text-2xl ml-2 icon-[pepicons-pop--arrow-right]") {}
        }
    )
}

#[component]
pub fn SectionOne<G: Html>(cx: Scope) -> View<G> {
    let title = "Bringing clarity to your Cosmos experience";
    let description = "Star Scope is a notification service for Cosmos that helps you stay updated on token transfers, new governance proposals, validator problems and ending of unbonding periods that are important for you.";

    view! {cx,
        div(class="hidden lg:flex justify-center flex-auto flex-shrink-0 min-h-screen bg-landing_page-bg") {
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
        div(class="lg:hidden flex flex-col items-center min-h-screen p-8 bg-landing_page-bg") {
            div(class="w-full") {
                div(class="flex justify-between items-center") {
                    img(class="h-auto w-48", src=keys::LOGO_WITH_TEXT_IMG) {}
                    LaunchButton {}
                }
            }
            div(class="flex flex-col grow items-center justify-center text-center") {
                div(class="flex flex-col items-center") {
                    h1(class="text-lg font-bold") {
                        (title)
                    }
                    p(class="text-sm my-8") {
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
pub fn SectionTwo<G: Html>(cx: Scope) -> View<G> {
    view! {cx,

    }
}

#[component]
pub fn Home<G: Html>(cx: Scope) -> View<G> {
    view! {cx,
        SectionOne {}
        // SectionTwo {}
    }
}
