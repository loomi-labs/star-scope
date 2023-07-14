use crate::components::messages::create_error_msg_from_status;
use crate::{AppState, Services};
use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;

#[component(inline_props)]
pub fn AskDeleteDialog<'a, G: Html>(cx: Scope<'a>, is_open: &'a Signal<bool>) -> View<G> {
    let app_state = use_context::<AppState>(cx);

    let handle_delete = |cx: Scope| {
        is_open.set(false);
        app_state.set_showing_dialog(false);
        spawn_local_scoped(cx, async move {
            let services = use_context::<Services>(cx);
            let request = services.grpc_client.create_request(());
            let response = services
                .grpc_client
                .get_user_service()
                .delete_account(request)
                .await;
            match response {
                Ok(_) => {
                    let app_state = use_context::<crate::AppState>(cx);
                    app_state.logout();
                }
                Err(status) => create_error_msg_from_status(cx, status),
            }
        });
    };

    create_effect(cx, move || {
        if *is_open.get() {
            app_state.set_showing_dialog(true); // sets the backdrop to be visible
        }
    });

    create_effect(cx, move || {
        if !(*app_state.is_dialog_open.get()) {
            is_open.set(false);
        }
    });

    view! {cx,
        dialog(class="fixed inset-0 bg-white p-4 rounded-lg z-40", open=*is_open.get()) {
            div(class="flex flex-col p-4") {
                div(class="flex flex-col items-center") {
                    span(class="w-12 h-12 text-black icon-[ph--trash]") {}
                    h2(class="text-lg font-semibold") { ("Delete account") }
                    p(class="my-4 text-center") { ("Are you sure you want to delete your account?") }
                }
                div(class="flex justify-center mt-2") {
                    button(class="border-2 border-gray-500 text-gray-500 hover:bg-gray-500 hover:text-white font-semibold px-4 py-2 rounded mr-2",
                            on:click=move |event: web_sys::Event| {
                        event.stop_propagation();
                        app_state.set_showing_dialog(false);
                    }) { "Cancel" }
                    button(class="bg-red-500 hover:bg-red-600 text-white font-semibold px-4 py-2 rounded",
                            on:click=move |_| handle_delete(cx)) { "Delete" }
                }
            }
        }
    }
}

#[component]
pub async fn Settings<G: Html>(cx: Scope<'_>) -> View<G> {
    let is_open = create_signal(cx, false);

    view! {cx,
        AskDeleteDialog(is_open=is_open)
        div(class="flex flex-col justify-center items-center") {
            button(class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded",
            on:click=move |_| is_open.set(true)) {
                "Delete Account"
            }
        }
    }
}
