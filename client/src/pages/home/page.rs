use sycamore::prelude::*;
use crate::config::keys;

#[component]
pub fn SectionOne<G: Html>(cx: Scope) -> View<G> {
    view! {cx,
        div(class="flex justify-center flex-auto flex-shrink-0 min-h-screen") {
            div(class="flex flex-col justify-between lg:max-w-screen-lg xl:max-w-screen-xl h-full w-full") {
                div(class="flex") {
                    div(class="flex flex-col justify-between w-1/3 min-h-screen") {
                        img(class="h-auto w-64 mt-8", src=keys::LOGO_WITH_TEXT_IMG) {}
                        div(class="flex flex-col justify-start pb-16") {
                            h1(class="text-5xl font-bold") {
                                "Bringing clarity to your Cosmos experience"
                            }
                            p(class="text-2xl my-8") {
                                "Stay updated on token transfers, new governance proposals, validator problems and ending of unbonding periods that are important for you."
                            }
                            button(class="rounded-xl flex items-center justify-center px-4 py-2 w-64 text-black hover:text-white bg-gradient-to-r from-primary_gradient-from to-primary_gradient-to") {
                                span(class="text-xl") {
                                    "Launch App"
                                }
                                span(class="text-2xl ml-2 icon-[pepicons-pop--arrow-right]") {}
                            }
                        }
                        div {}
                    }
                    div(class="flex flex-col justify-center w-2/3") {
                        img(class="h-full w-full object-contain", src=keys::SCOPE_IMG) {}
                    }
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
        SectionTwo {}
    }
}
