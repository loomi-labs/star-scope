use gloo_timers::future::TimeoutFuture;
use log::error;
use sycamore::futures::spawn_local_scoped;
use sycamore::motion::create_raf;
use sycamore::prelude::*;
use tonic::Status;

use crate::{AppState, InfoLevel, InfoMsg};

#[component(inline_props)]
pub fn Message<G: Html>(cx: Scope, msg: InfoMsg, style: String, timeout: u32) -> View<G> {
    let app_state = use_context::<AppState>(cx);

    let state = create_signal(cx, 1.0);
    let start_elapse = timeout as f64 * 0.7 * 1000.0; // start elapsing after 70% of the timeout
    let elapse_time = timeout as f64 * 1000.0 - start_elapse;
    let (_running, start, stop) = create_raf(cx, move || {
        let elapsed = js_sys::Date::now() - msg.created_at;
        if elapsed > start_elapse {
            let new_state = 1.0 - (elapsed - start_elapse) / elapse_time;
            state.set(new_state);
            if new_state <= 0.0 {
                app_state.remove_message(msg.id);
            }
        }
    });
    start();

    let color = match msg.level {
        InfoLevel::Info => "bg-msg-blue border-msg-blue text-msg-blue",
        InfoLevel::Success => "bg-msg-green border-msg-green text-msg-green",
        InfoLevel::Error => "bg-msg-red border-msg-red text-msg-red",
    };

    let hover_color = match msg.level {
        InfoLevel::Info => "hover:bg-msg-blue",
        InfoLevel::Success => "hover:bg-msg-green",
        InfoLevel::Error => "hover:bg-msg-red",
    };

    let icon = match msg.level {
        InfoLevel::Info => "icon-[material-symbols--error-rounded]",
        InfoLevel::Success => "icon-[ic--round-check-circle]",
        InfoLevel::Error => "icon-[ph--x-circle-fill]",
    };

    view! { cx,
        div(class=format!("absolute bottom-0 right-0 flex items-center justify-center min-w-96 p-4 m-6 rounded-lg bg-white border-l-[20px] drop-shadow-lg {}", color), style=format!("{} opacity: {}", style, state.get().as_ref())) {
            span(class=format!("w-10 h-10 {}", icon)) {}
            div(class="flex flex-col pl-4") {
                h3(class="text-lg font-bold") { (msg.title) }
                p(class="text-sm text-black") { (msg.message) }
            }
            button(
                class="top-0 right-0 self-start pointer-events-auto",
                on:click=move |_| {
                    spawn_local_scoped(cx, async move {
                        stop();
                        app_state.remove_message(msg.id);
                    });
                }
            ) {
                div(class=format!("flex items-center justify-center rounded-full text-white w-8 h-8 {}", hover_color)) {
                    span(class="w-6 h-6 bg-black icon-[bi--x]") {}
                }
            }
        }
    }
}

#[derive(Debug, Clone, PartialEq)]
pub struct IndexedItem<T> {
    index: usize,
    item: T,
}

#[component]
pub fn MessageOverlay<G: Html>(cx: Scope) -> View<G> {
    let app_state = use_context::<AppState>(cx);
    let messages = create_selector(cx, || {
        app_state
            .messages
            .get()
            .iter()
            .enumerate()
            .map(|(index, msg)| IndexedItem {
                index,
                item: msg.clone(),
            })
            .collect::<Vec<_>>()
    });

    view!(
        cx,
        div(class="fixed inset-0 min-h-[100svh] flex justify-center items-center flex-auto flex-shrink-0 z-50 pointer-events-none") {
            div(class="relative flex flex-col lg:max-w-screen-lg xl:max-w-screen-xl h-full w-full") {
                Indexed(
                    iterable = messages,
                    view = move |cx, iItem| {
                        let style = format!("margin-bottom: {}rem;", 7 * iItem.index + 2);
                        let item = iItem.item.get();

                        view!{cx,
                            Message(
                                msg=item.as_ref().clone(),
                                style=style,
                                timeout=item.timeout,
                            )
                        }
                    },
                )
            }
        }
    )
}

pub fn create_timed_message(
    cx: Scope,
    title: impl Into<String>,
    message: impl Into<String>,
    level: InfoLevel,
    timeout: u32,
) {
    let app_state = use_context::<AppState>(cx);
    let title = title.into();
    let message = message.into();
    let uuid = app_state.add_message(title.clone(), message.clone(), level.clone(), timeout);
    if level == InfoLevel::Error {
        error!("{}: {}", title, message);
    }
    create_effect(cx, move || {
        spawn_local_scoped(cx, async move {
            TimeoutFuture::new(1000 * timeout).await;
            app_state.remove_message(uuid);
        });
    });
}

pub fn create_message(
    cx: Scope,
    title: impl Into<String>,
    message: impl Into<String>,
    level: InfoLevel,
) {
    create_timed_message(cx, title, message, level, 10)
}

pub fn create_error_msg_from_status(cx: Scope, status: Status) {
    create_message(
        cx,
        status.code().to_string().as_str(),
        status.message(),
        InfoLevel::Error,
    );
}
