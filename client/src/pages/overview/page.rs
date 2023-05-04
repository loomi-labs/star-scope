use sycamore::prelude::*;

use crate::components::messages::create_error_msg_from_status;
use crate::Services;

#[component]
pub async fn Overview<G: Html>(cx: Scope<'_>) -> View<G> {
    // provide_context(cx, OverviewState::new());
    //
    // query_chains(cx.to_owned()).await;
    //
    // view! {cx,
    //     div(class="flex flex-col h-full w-full p-8") {
    //         h1(class="text-2xl font-bold") { "Overview" }
    //         div(class="flex flex-col p-8 bg-white dark:bg-gray-600 rounded-lg shadow") {
    //             Search {}
    //         }
    //     }
    // }
    view!(cx,
    "Not implemented yet"
    )
}
