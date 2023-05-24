use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;
use tonic::{Status};
use crate::{Services};
use crate::components::messages::create_error_msg_from_status;

#[component]
pub async fn Settings<G: Html>(cx: Scope<'_>) -> View<G> {
    let handle_delete = |cx: Scope| {
        spawn_local_scoped(cx, async move {
            let services = use_context::<Services>(cx);
            let request = services.grpc_client.create_request(());
            let response = services.grpc_client.
                get_user_service().
                delete_account(request).await;
            match response {
                Ok(_) => {
                    let app_state = use_context::<crate::AppState>(cx);
                    app_state.logout();
                },
                Err(status) => create_error_msg_from_status(cx, status),
            }
        });
    };

    view! {cx,
        div(class="flex flex-col justify-center items-center") {
            button(class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded",
            on:click=move |_| handle_delete(cx)) {
                "Delete Account"
            }
        }
    }
}
