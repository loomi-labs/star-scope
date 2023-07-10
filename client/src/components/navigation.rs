use sycamore::prelude::*;

use crate::config::keys;
use crate::pages::notifications::page::NotificationsState;
use crate::types::protobuf::event::EventType;
use crate::utils::url::safe_navigate;
use crate::{AppRoutes, AppState, EventsState};

#[component]
pub fn Header<G: Html>(cx: Scope) -> View<G> {
    let app_state = use_context::<AppState>(cx);
    view!(cx,
        div(class="flex items-center justify-between h-14 text-white dark:text-purple-600") {
            div(class="flex flex-grow items-center justify-start pl-3 h-14") {
                button(on:click=move |_| safe_navigate(cx, AppRoutes::Home), class="relative") {
                        img(src=keys::LOGO_WITH_TEXT_WHITE_IMG, class="h-auto w-44 transition-transform duration-300 transform")
                        img(src=keys::LOGO_WITH_TEXT_ORANGE_IMG, class="h-auto w-44 absolute top-0 left-0 opacity-0 transition-opacity duration-300 transform hover:opacity-100")
                }
            }
            div(class="flex justify-between items-center h-14 pr-8") {
                ul(class="flex items-center") {
                    li {
                        button(class="flex items-center p-2 rounded hover:text-primary dark:hover:text-primary", on:click=move |_| app_state.logout()) {
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

    let handle_notification_click = |cx: Scope, event_type: Option<EventType>| {
        notifications_state.apply_filter(event_type);
        safe_navigate(cx, AppRoutes::Notifications);
    };

    fn get_event_count(events_state: &EventsState, event_type: EventType) -> Option<i32> {
        events_state
            .event_count_map
            .get()
            .get(&event_type)
            .map(|e| e.unread_count)
            .filter(|e| *e > 0)
    }

    let cnt_funding = create_selector(cx, move || {
        get_event_count(events_state, EventType::Funding)
    });
    let cnt_staking = create_selector(cx, move || {
        get_event_count(events_state, EventType::Staking)
    });
    let cnt_dex = create_selector(cx, move || get_event_count(events_state, EventType::Dex));
    let cnt_governance = create_selector(cx, move || {
        get_event_count(events_state, EventType::Governance)
    });
    let cnt_all = create_selector(cx, || {
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

    fn is_active_notification_route(
        event_type: Option<EventType>,
        notifications_state: &NotificationsState,
        active_route: Option<AppRoutes>,
    ) -> bool {
        active_route == Some(AppRoutes::Notifications)
            && notifications_state.has_filter_applied(event_type)
    }

    let button_class = "relative flex flex-row items-center text-center max-w-full h-11 pr-6";
    let button_interactivity_class = "focus:outline-none hover:bg-blue-800 dark:hover:bg-purple-800 dark:hover:text-primary text-white-600 hover:text-white-800";
    let span_icon_class = "inline-flex justify-center items-center ml-4";
    let span_text_class = "overflow-y-auto overflow-x-hidden ml-2 text-base tracking-wide truncate";
    let badge_class = "inline-flex items-center justify-center w-5 h-5 ml-0 rounded-full text-[10px] font-bold text-white bg-red-500 border-2 border-white dark:border-gray-900";

    let n_button_data = vec![
        (None, "icon-[lucide--copy-check]", "All", cnt_all),
        (
            Some(EventType::Funding),
            "icon-[ep--coin]",
            "Funding",
            cnt_funding,
        ),
        (
            Some(EventType::Staking),
            "icon-[carbon--equalizer]",
            "Staking",
            cnt_staking,
        ),
        (
            Some(EventType::Dex),
            "icon-[fluent--money-24-regular]",
            "Dex",
            cnt_dex,
        ),
        (
            Some(EventType::Governance),
            "icon-[icon-park-outline--palace]",
            "Governance",
            cnt_governance,
        ),
    ];

    let notification_button_views = View::new_fragment(
        n_button_data.iter().map(|&d| view! { cx, li {
            button(on:click=move |_| handle_notification_click(cx, d.0), class=format!("{} {} {}", button_class, button_interactivity_class, if is_active_notification_route(d.0, notifications_state, *app_state.route.get()) { "text-primary" } else { "" })) {
                span(class=format!("{} w-1 h-5 rounded-r-lg absolute", if is_active_notification_route(d.0, notifications_state, *app_state.route.get()) {"bg-primary"} else { "" })) {}
                span(class=format!("{} {} {}", d.1, span_icon_class, if *is_sidebar_hovered.get() { "ml-2" } else { "ml-4" })) {
                    div(class="w-16 h-16")
                }
                span(class=span_text_class) { (d.2) }
                (if d.3.get().is_some() {
                    view! {cx,
                        div(class="absolute top-0 right-1") {
                            div(class=badge_class) { (d.3.get().unwrap_or(0)) }
                        }
                    }
                    } else {
                    view! {cx, div()}
                })
            }
        } }).collect()
    );

    let s_button_data = vec![
        (
            Some(AppRoutes::Communication),
            "icon-[mi--message]",
            "Communication",
        ),
        (
            Some(AppRoutes::Settings),
            "icon-[iconamoon--profile]",
            "Account",
        ),
    ];

    let settings_button_views = View::new_fragment(
        s_button_data.iter().map(|&d| {
            view! { cx,
                li {
                    button(
                        on:click=move |_| safe_navigate(cx, d.0.unwrap()),
                        class=format!("{} {} {}", button_class, button_interactivity_class, if *app_state.route.get() == d.0 { "text-primary" } else { "" })
                    ) {
                        span(class=format!("{} w-1 h-5 rounded-r-lg absolute", if *app_state.route.get() == d.0 {"bg-primary"} else { "" })) {}
                        span(class=format!("{} {} {}", d.1, span_icon_class, if *is_sidebar_hovered.get() { "ml-2" } else { "ml-4" })) {
                            div(class="w-16 h-16") {}
                        }
                        span(class=span_text_class) { (d.2) }
                    }
                }
            }
        }).collect()
    );

    view! { cx,
        div(class="flex flex-col justify-between w-14 hover:w-64 lg:w-64 text-white transition-all duration-300 border-none",
            on:mouseenter=move |_| is_sidebar_hovered.set(true),
            on:mouseleave=move |_| is_sidebar_hovered.set(false),
        ) {
            ul(class="py-4 space-y-1 dark:bg-purple-800 rounded-b-lg") {
                li() {
                    a(href=AppRoutes::Notifications, class=format!("{} {} transition duration-500 ease-in-out text-purple-600 lg:text-purple-600", button_class, if *is_sidebar_hovered.get() { "" } else { "text-purple-600/0" })) {
                        div(style="overflow: hidden; text-overflow: ellipsis;") {
                            span(class=format!("ml-4 text-base tracking-wide")) { "Notifications" }
                        }
                    }
                    ul() {
                        (notification_button_views)
                    }
                }
            }
            ul(class="pt-4 pb-10 md:pb-4 space-y-1 dark:bg-purple-800 rounded-t-lg") {
                li() {
                    a(href=AppRoutes::Notifications, class=format!("{} {} transition duration-500 ease-in-out text-purple-600 lg:text-purple-600", button_class, if *is_sidebar_hovered.get() { "" } else { "text-purple-600/0" })) {
                        div(style="overflow: hidden; text-overflow: ellipsis;") {
                            span(class=format!("ml-4 text-base tracking-wide")) { "Settings" }
                        }
                    }
                    ul() {
                        (settings_button_views)
                    }
                }
            }
        }
    }
}

#[component(inline_props)]
pub fn Navigation<'a, G: Html>(cx: Scope<'a>, children: Children<'a, G>) -> View<G> {
    let children = children.call(cx);
    view! { cx,
        div(class="min-h-[100svh] max-h-[100lvh] flex justify-center items-center flex-auto flex-shrink-0") {
            div(class="flex flex-col lg:max-w-screen-lg xl:max-w-screen-xl h-full w-full dark:bg-purple-800") {
                Header{}
                div(class="flex flex-row h-full w-full dark:bg-d-bg") {
                    Sidebar{}
                    div(class="w-full p-4 md:p-8 lg:p-0 lg:py-8 lg:pl-8 md:max-w-auto h-[calc(100vh-theme(space.16))] overflow-y-auto overflow-x-visible") {    // TODO: fix the 100vh-theme(space.16) hack
                        (children)
                    }
                }
            }
        }
    }
}
