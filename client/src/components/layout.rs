use sycamore::prelude::*;

use crate::{AppRoutes, AppState, EventsState};
use crate::pages::notifications::page::NotificationsState;
use crate::types::types::grpc::EventType;
use crate::utils::url::safe_navigate;

#[component]
pub fn Header<G: Html>(cx: Scope) -> View<G> {
    let app_state = use_context::<AppState>(cx);
    let logo_color = "#D68940";
    view!(cx,
        div(class="flex items-center justify-between h-14 mr-8 z-10 text-white dark:text-purple-600") {
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
                div(class="outline-none focus:outline-none") {}
                div(class="w-full pl-3 text-sm text-black outline-none focus:outline-none bg-transparent" ) {}
            }
            div(class="flex justify-between items-center h-14 header-right") {
                ul(class="flex items-center") {
                    li {
                        button(class="flex items-center mr-4 p-2 rounded hover:text-primary dark:hover:text-primary", on:click=move |_| app_state.logout()) {
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

#[component]
pub fn Sidebar<G: Html>(cx: Scope) -> View<G> {
    let notifications_state = use_context::<NotificationsState>(cx);
    let app_state = use_context::<AppState>(cx);
    let events_state = use_context::<EventsState>(cx);

    let handle_click = |cx: Scope, event_type: Option<EventType>| {
        notifications_state.apply_filter(event_type);
        safe_navigate(cx, AppRoutes::Notifications);
    };

    fn get_event_count(events_state: &EventsState, event_type: EventType) -> Option<i32> {
        events_state
            .event_count_map
            .get()
            .get(&event_type)
            .map(|e| e.unread_count.clone())
            .filter(|e| *e > 0)
    }

    let cnt_funding = create_memo(cx, move || get_event_count(&events_state, EventType::Funding));
    let cnt_staking = create_memo(cx, move || get_event_count(&events_state, EventType::Staking));
    let cnt_dex = create_memo(cx, move || get_event_count(&events_state, EventType::Dex));
    let cnt_governance = create_memo(cx, move || get_event_count(&events_state, EventType::Governance));
    let cnt_all = create_memo(cx, || {
        let all = cnt_funding.get().as_ref().unwrap_or_else(|| 0)
            + cnt_staking.get().as_ref().unwrap_or_else(|| 0)
            + cnt_dex.get().as_ref().unwrap_or_else(|| 0)
            + cnt_governance.get().as_ref().unwrap_or_else(|| 0);
        if all > 0 {
            Some(all)
        } else {
            None
        }
    });

    let is_sidebar_hovered = create_signal(cx, false);

    fn is_active_notification_route(event_type: Option<EventType>, notifications_state: &NotificationsState, active_route: &AppRoutes) -> bool {
        active_route == &AppRoutes::Notifications && notifications_state.has_filter_applied(event_type)
    }

    let button_class = "relative flex flex-row items-center text-center max-w-full h-11 pr-6";
    let button_interactivity_class = "focus:outline-none hover:bg-blue-800 dark:hover:bg-purple-800 dark:hover:text-primary text-white-600 hover:text-white-800";
    let span_icon_class = "inline-flex justify-center items-center ml-4";
    let span_text_class = "overflow-y-auto overflow-x-hidden ml-2 text-sm tracking-wide truncate";
    let badge_class = "inline-flex items-center justify-center w-5 h-5 ml-0 rounded-full text-[10px] font-bold text-white bg-red-500 border-2 border-white dark:border-gray-900";

    let button_data = vec![
        (None, "icon-[lucide--copy-check]", "All", cnt_all),
        (Some(EventType::Funding), "icon-[ep--coin]", "Funding", cnt_funding),
        (Some(EventType::Staking), "icon-[carbon--equalizer]", "Staking", cnt_staking),
        (Some(EventType::Dex), "icon-[fluent--money-24-regular]", "Dex", cnt_dex),
        (Some(EventType::Governance), "icon-[icon-park-outline--palace]", "Governance", cnt_governance),
    ];

    let notification_button_views = View::new_fragment(
        button_data.iter().map(|&d| view! { cx, li {
            button(on:click=move |_| handle_click(cx, d.0), class=format!("{} {} {}", button_class, button_interactivity_class, if is_active_notification_route(d.0, notifications_state, app_state.route.get().as_ref()) { "text-primary" } else { "" })) {
                span(class=format!("{} w-1 h-5 rounded-r-lg absolute", if is_active_notification_route(d.0, notifications_state, app_state.route.get().as_ref()) {"bg-primary"} else { "" })) {}
                span(class=format!("{} {} {}", d.1, span_icon_class, if *is_sidebar_hovered.get() { "ml-2" } else { "ml-4" })) {
                    div(class="w-16 h-16")
                }
                span(class=span_text_class) { (d.2) }
                (if d.3.get().is_some() {
                    view! {cx,
                        div(class="absolute top-0 right-1") {
                            div(class=badge_class) { (d.3.get().unwrap()) }
                        }
                    }
                    } else {
                    view! {cx, div()}
                })
            }
        } }).collect()
    );

    view! { cx,
        div(class="h-full flex flex-col top-14 left-0 w-14 hover:w-64 lg:w-64 text-white transition-all duration-300 border-none z-10",
            on:mouseenter=move |_| is_sidebar_hovered.set(true),
            on:mouseleave=move |_| is_sidebar_hovered.set(false),
        ) {
            div(class="flex flex-col") {
                ul(class="flex flex-col py-4 space-y-1 dark:bg-purple-800 rounded") {
                    li() {
                        a(href=AppRoutes::Notifications, class=format!("{} {} transition duration-500 ease-in-out text-purple-600 lg:text-purple-600", button_class, if *is_sidebar_hovered.get() { "" } else { "text-purple-600/0" })) {
                            div(style="overflow: hidden; text-overflow: ellipsis;") {
                                span(class=format!("ml-3 text-base tracking-wide")) { "Notifications" }
                            }
                        }
                        ul() {
                            (notification_button_views)
                        }
                    }
                }
            }
            div(class="flex flex-col", style="height: calc(100vh - 460px)")
            div(class="flex flex-col pb-10") {
                ul(class="flex flex-col py-2 space-y-1 dark:bg-purple-800 rounded") {
                    li() {
                        a(href=AppRoutes::Settings, class=format!("{} {} {}", button_class, button_interactivity_class, if *app_state.route.get() == AppRoutes::Settings {"text-primary"} else { "" })) {
                            span(class=format!("{} w-1 h-5 rounded-r-lg absolute", if *app_state.route.get() == AppRoutes::Settings {"bg-primary"} else { "" })) {}
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
                    div(class="p-8 w-full max-w-[90vw] md:max-w-auto h-[calc(100vh-theme(space.16))] overflow-y-auto overflow-x-visible") {
                        (children)
                    }
                }
            }
        }
    }
}
