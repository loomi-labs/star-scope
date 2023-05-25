use sycamore::prelude::*;

use crate::{AppRoutes, AppState};
use crate::pages::notifications::page::NotificationsState;
use crate::types::types::grpc::EventType;
use crate::utils::url::safe_navigate;

#[component]
pub fn Header<G: Html>(cx: Scope) -> View<G> {
    let app_state = use_context::<AppState>(cx);
    let logo_color = "#D68940";
    view!(cx,
        div(class="flex items-center justify-between h-14 mr-8 text-white z-10 dark:bg-purple-800") {
            div(class="flex items-center justify-start pl-4 h-14 w-14 lg:w-64 border-none") {
                svg(xmlns="http://www.w3.org/2000/svg", width="60", height="60", viewBox="0 0 563 547.33") {
                    path(fill=logo_color, d="M282.22,24c3.39,42.11,5.12,84.23,6.35,126.34s1.69,84.22,1.64,126.33-.6,84.22-1.77,126.33-2.95,84.23-6.22,126.34c-3.26-42.11-5-84.23-6.22-126.34s-1.73-84.22-1.76-126.33.45-84.22,1.64-126.33S278.84,66.07,282.22,24Z")
                    path(fill=logo_color, d="M538.55,275.5c-42.11,3.35-84.22,5.05-126.33,6.24s-84.23,1.63-126.34,1.54-84.22-.66-126.33-1.87-84.22-3-126.33-6.32c42.12-3.23,84.23-4.9,126.34-6.11s84.22-1.67,126.33-1.67,84.23.53,126.34,1.74S496.45,272.08,538.55,275.5Z")
                    path(fill=logo_color, d="M338.13,221.33c-7.79,10.41-16.23,20.16-24.87,29.71s-17.58,18.8-26.73,27.84S268,296.76,258.45,305.35,239,322.29,228.6,330.05c7.84-10.37,16.29-20.11,24.92-29.67s17.57-18.81,26.68-27.88,18.44-17.94,28-26.52S327.66,229,338.13,221.33Z")
                    path(fill=logo_color, d="M338.76,331c-10.41-7.79-20.15-16.24-29.69-24.89s-18.79-17.6-27.82-26.75-17.86-18.51-26.45-28.11-16.92-19.43-24.67-29.86c10.35,7.84,20.09,16.3,29.64,24.94s18.8,17.58,27.86,26.7,17.93,18.45,26.51,28.06S331.06,320.53,338.76,331Z")
                    path(fill="none", stroke=logo_color, stroke-miterlimit="10", stroke-width="30", d="M138.39,272.21A143.53,143.53,0,0,1,281.91,130.89")
                    path(fill="none", stroke=logo_color, stroke-miterlimit="10", stroke-width="10", d="M284.89,417.93c-1,0-2,0-3,0A143.53,143.53,0,0,1,138.37,274.43")
                    path(fill="none", stroke=logo_color, stroke-miterlimit="10", stroke-width="30", d="M425.44,274.43a143.53,143.53,0,0,1-140.55,143.5")
                    path(fill="none", stroke=logo_color, stroke-miterlimit="10", stroke-width="10", d="M284.89,130.92A143.54,143.54,0,0,1,425.43,272.21")
                }
            }
            div(class="flex flex-grow justify-between items-center h-14 header-right") {
                // span(class="hidden lg:block") { (app_state.get_user_name()) }
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

fn highlight_active_notification_route(event_type: Option<EventType>, notifications_state: &NotificationsState, active_route: &AppRoutes) -> String {
    if active_route == &AppRoutes::Notifications {
        if notifications_state.has_filter_applied(event_type) {
            return "text-primary".to_string();
        }
    }
    "".to_string()
}

fn highlight_active_route(active_route: &AppRoutes, route: &AppRoutes) -> String {
    if active_route == route {
        return "text-primary".to_string();
    }
    "".to_string()
}

#[component]
pub fn Sidebar<G: Html>(cx: Scope) -> View<G> {
    let notifications_state = use_context::<NotificationsState>(cx);
    let app_state = use_context::<AppState>(cx);

    let button_class = "relative flex flex-row items-center max-w-full h-11 focus:outline-none hover:bg-blue-800 dark:hover:bg-purple-800 dark:hover:text-primary text-white-600 hover:text-white-800 border-l-4 border-transparent pr-6";
    let span_icon_class = "inline-flex justify-center items-center ml-4 font-size-20";
    let span_text_class = "ml-2 text-sm tracking-wide truncate";

    let handle_click = |cx: Scope, event_type: Option<EventType>| {
        notifications_state.apply_filter(event_type);
        safe_navigate(cx, AppRoutes::Notifications);
    };
    view! { cx,
        div(class="h-full flex flex-col top-14 left-0 w-14 hover:w-64 lg:w-64 h-full text-white transition-all duration-300 border-none z-10") {
            div(class="overflow-y-auto overflow-x-hidden flex flex-col") {
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
                                button(on:click=move |_| handle_click(cx, None), class=format!("{} {}", button_class, highlight_active_notification_route(None, notifications_state, app_state.route.get().as_ref()))) {
                                    span(class=format!("{} icon-[lucide--copy-check]", span_icon_class)) {
                                        div(class="w-16 h-16")
                                    }
                                    span(class=span_text_class) { "All" }
                                }
                            }
                            li() {
                                button(on:click=move |_| handle_click(cx, Some(EventType::Funding)), class=format!("{} {}", button_class, highlight_active_notification_route(Some(EventType::Funding), notifications_state, app_state.route.get().as_ref()))) {
                                    span(class=format!("{} icon-[ep--coin]", span_icon_class)) {
                                        div(class="w-16 h-16")
                                    }
                                    span(class=span_text_class) { "Funding" }
                                }
                            }
                            li() {
                                button(on:click=move |_| handle_click(cx, Some(EventType::Staking)), class=format!("{} {}", button_class, highlight_active_notification_route(Some(EventType::Staking), notifications_state, app_state.route.get().as_ref()))) {
                                    span(class=format!("{} icon-[arcticons--coinstats]", span_icon_class)) {
                                        div(class="w-16 h-16")
                                    }
                                    span(class=span_text_class) { "Staking" }
                                }
                            }
                            li() {
                                button(on:click=move |_| handle_click(cx, Some(EventType::Dex)), class=format!("{} {}", button_class, highlight_active_notification_route(Some(EventType::Dex), notifications_state, app_state.route.get().as_ref()))) {
                                    span(class=format!("{} icon-[fluent--money-24-regular]", span_icon_class)) {
                                        div(class="w-16 h-16")
                                    }
                                    span(class=span_text_class) { "DEX'es" }
                                }
                            }
                            li() {
                                button(on:click=move |_| handle_click(cx, Some(EventType::Governance)), class=format!("{} {}", button_class, highlight_active_notification_route(Some(EventType::Governance), notifications_state, app_state.route.get().as_ref()))) {
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
            div(class="flex flex-col", style="height: calc(100vh - 460px)")
            div(class="overflow-y-auto overflow-x-hidden flex flex-col pb-10") {
                ul(class="flex flex-col py-2 space-y-1 dark:bg-purple-800 rounded") {
                    li() {
                        a(href=AppRoutes::Settings, class=format!("{} {}", button_class, highlight_active_route(&app_state.route.get().as_ref(), &AppRoutes::Settings))) {
                            span(class=format!("{} icon-[streamline--interface-setting-cog-work-loading-cog-gear-settings-machine]", span_icon_class)) {
                                div(class="w-16 h-16")
                            }
                            span(class=span_text_class) { "Settings" }
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
                    div(class="p-8 w-full max-w-[90vw] md:max-w-auto h-[calc(100vh-theme(space.16))] overflow-y-auto") {
                        (children)
                    }
                }
            }
        }
    }
}
