use sycamore::futures::spawn_local_scoped;
use sycamore::prelude::*;

#[component(inline_props)]
pub fn Search<'a, G: Html>(
    cx: Scope<'a>,
    search_term: &'a Signal<String>,
    has_input_focus: &'a Signal<bool>,
    placeholder: &'a str,
) -> View<G> {
    view! {cx,
        input(
            class="w-full placeholder:italic placeholder:text-slate-400 block border border-slate-300 rounded-full px-4 py-2
                shadow-sm focus:outline-none focus:border-primary focus:ring-primary focus:ring-1 sm:text-sm",
            placeholder=placeholder,
            type="text",
            bind:value=search_term,
            on:focusin=move |_| has_input_focus.set(true),
            on:blur=move |_| {
                // this has to be delayed because otherwise the blur event will fire before the focusin event on the dialog
                spawn_local_scoped(cx, async move {
                    gloo_timers::future::TimeoutFuture::new(100).await;
                    has_input_focus.set(false)
                });
            },
        )
    }
}