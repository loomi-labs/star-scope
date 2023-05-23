use sycamore::prelude::*;

use crate::{AppRoutes, AppState};
use crate::pages::notifications::page::NotificationsState;
use crate::services::grpc::EventType;
use crate::utils::url::safe_navigate;

#[component]
pub fn Header<G: Html>(cx: Scope) -> View<G> {
    let app_state = use_context::<AppState>(cx);

    view!(cx,
        div(class="flex items-center justify-between h-14 mr-8 text-white z-10 dark:bg-purple-800") {
            div(class="flex items-center justify-start pl-4 h-14 w-14 lg:w-64 border-none") {
                span(class="icon-[game-icons--ringed-planet] h-10 w-10 text-primary") {}
            }
            div(class="flex flex-grow justify-between items-center h-14 header-right") {
                span(class="hidden lg:block") { (app_state.get_user_name()) }
                div(class="outline-none focus:outline-none") {}
                div(class="w-full pl-3 text-sm text-black outline-none focus:outline-none bg-transparent" ) {}
            }
            div(class="flex justify-between items-center h-14 header-right") {
                ul(class="flex items-center") {
                    li {
                        button(aria-hidde="true", class="group w-9 h-9 transition-colors duration-200 rounded-full shadow-md bg-blue-200 hover:bg-primary dark:bg-purple-700 dark:hover:text-primary focus:outline-none") {
                            i(class="fas fa-bell text-lg") {}
                        }
                    }
                    li {
                        div(class="block w-px h-6 mx-3 bg-gray-400 dark:bg-purple-600") {}
                    }
                    li {
                        button(class="flex items-center mr-4 p-2 rounded hover:text-primary dark:hover:bg-purple-700 dark:hover:text-primary", on:click=move |_| app_state.logout()) {
                            span(class="inline-flex mr-1") {
                                i(class="fas fa-sign-out-alt text-xl") {}
                            }
                            "Logout"
                        }
                    }
                }
            }
        }
    )
}

fn highlight_active_route(event_type: Option<EventType>, notifications_state: &NotificationsState) -> String {
    if notifications_state.has_filter_applied(event_type) {
        "text-primary".to_string()
    } else {
        "".to_string()
    }
}

#[component]
pub fn Sidebar<G: Html>(cx: Scope) -> View<G> {
    let notifications_state = use_context::<NotificationsState>(cx);

    let button_class = "relative flex flex-row items-center max-w-full h-11 focus:outline-none hover:bg-blue-800 dark:hover:bg-purple-800 dark:hover:text-primary text-white-600 hover:text-white-800 border-l-4 border-transparent pr-6";
    let span_icon_class = "inline-flex justify-center items-center ml-4 font-size-20";
    let span_text_class = "ml-2 text-sm tracking-wide truncate";

    let handle_click = |cx: Scope, event_type: Option<EventType>| {
        notifications_state.apply_filter(event_type);
        safe_navigate(cx, AppRoutes::Notifications);
    };

    view! { cx,
        div(class="h-full flex flex-col top-14 left-0 w-14 hover:w-64 lg:w-64 h-full text-white transition-all duration-300 border-none z-10 sidebar") {
            div(class="overflow-y-auto overflow-x-hidden flex flex-col justify-between flex-grow") {
                ul(class="flex flex-col py-4 space-y-1 dark:bg-purple-800 rounded") {
                    li() {
                        a(href=AppRoutes::Notifications, class=button_class) {
                            span(class=format!("{} icon-[mdi--bell]", span_icon_class)) {
                                div(class="w-16 h-16")
                            }
                            span(class=format!("{} uppercase", span_text_class)) { "Notifications" }
                        }
                        ul() {
                            li() {
                                button(on:click=move |_| handle_click(cx, None), class=format!("{} {}", button_class, highlight_active_route(None, notifications_state))) {
                                    span(class=format!("{} icon-[lucide--copy-check]", span_icon_class)) {
                                        div(class="w-16 h-16")
                                    }
                                    span(class=span_text_class) { "All" }
                                }
                            }
                            li() {
                                button(on:click=move |_| handle_click(cx, Some(EventType::Funding)), class=format!("{} {}", button_class, highlight_active_route(Some(EventType::Funding), notifications_state))) {
                                    span(class=format!("{} icon-[ep--coin]", span_icon_class)) {
                                        div(class="w-16 h-16")
                                    }
                                    span(class=span_text_class) { "Funding" }
                                }
                            }
                            li() {
                                button(on:click=move |_| handle_click(cx, Some(EventType::Staking)), class=format!("{} {}", button_class, highlight_active_route(Some(EventType::Staking), notifications_state))) {
                                    span(class=format!("{} icon-[arcticons--coinstats]", span_icon_class)) {
                                        div(class="w-16 h-16")
                                    }
                                    span(class=span_text_class) { "Staking" }
                                }
                            }
                            li() {
                                button(on:click=move |_| handle_click(cx, Some(EventType::Dex)), class=format!("{} {}", button_class, highlight_active_route(Some(EventType::Dex), notifications_state))) {
                                    span(class=format!("{} icon-[fluent--money-24-regular]", span_icon_class)) {
                                        div(class="w-16 h-16")
                                    }
                                    span(class=span_text_class) { "DEX'es" }
                                }
                            }
                            li() {
                                button(on:click=move |_| handle_click(cx, Some(EventType::Governance)), class=format!("{} {}", button_class, highlight_active_route(Some(EventType::Governance), notifications_state))) {
                                    span(class=format!("{} icon-[icon-park-outline--palace]", span_icon_class)) {
                                        div(class="w-16 h-16")
                                    }
                                    span(class=span_text_class) { "Governance" }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}

#[component(inline_props)]
pub fn LayoutWrapper<'a, G: Html>(cx: Scope<'a>, children: Children<'a, G>) -> View<G> {
    let children = children.call(cx);
    view! { cx,
        div(class="min-h-screen flex justify-center items-center flex-auto flex-shrink-0 antialiased dark:bg-purple-900") {
            div(class="flex flex-col lg:max-w-screen-lg xl:max-w-screen-xl h-full w-full") {
                Header{}
                div(class="flex flex-row h-full w-full") {
                    Sidebar{}
                    div(class="p-8 max-w-[90vw] md:max-w-auto h-full w-full") {
                        (children)
                    }
                }
            }
        }
    }
}
