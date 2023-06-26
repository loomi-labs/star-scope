use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;

#[component]
pub fn Setup<G: Html>(cx: Scope) -> View<G> {
    view! {cx,
        div()
    }
}
