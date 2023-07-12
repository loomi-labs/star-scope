use std::collections::HashSet;
use std::fmt::Display;
use std::hash::Hash;

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
        div(class="relative flex items-center") {
            input(
                class="relative w-full placeholder:italic placeholder:text-slate-400 block border border-slate-300 rounded-full px-4 py-2
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
            span(class="absolute right-3 w-5 h-5 bg-slate-400 icon-[ion--search] pointer-events-none")
        }
    }
}

#[derive(Debug, PartialEq, Clone)]
pub struct Searchable<'a, T> {
    pub entity: T,
    pub is_selected: &'a Signal<bool>,
}

#[derive(Prop)]
pub struct SearchEntityProps<'a, T> {
    searchables: Vec<Searchable<'a, T>>,
    selected_entities: &'a Signal<HashSet<T>>,
    placeholder: &'a str,
    #[builder(default = false)]
    show_results_for_empty_search: bool,
}

#[component]
pub fn SearchEntity<'a, G: Html, T>(cx: Scope<'a>, props: SearchEntityProps<'a, T>) -> View<G>
where
    T: Clone + Hash + Eq + PartialEq + Display + 'a,
{
    const MAX_ENTRIES: usize = 10;
    let search_term = create_signal(cx, String::new());

    let searchables_ref = create_ref(cx, props.searchables);
    let search_results = create_selector(cx, move || {
        let search = search_term.get().to_lowercase();
        let mut results = vec![];
        if !search.is_empty() {
            for searhable in searchables_ref.iter() {
                if searhable
                    .entity
                    .to_string()
                    .to_ascii_lowercase()
                    .contains(&search)
                {
                    results.push(searhable.clone());
                }
                if results.len() >= MAX_ENTRIES {
                    break;
                }
            }
        } else if props.show_results_for_empty_search {
            results = searchables_ref.iter().take(MAX_ENTRIES).cloned().collect();
        }
        results
    });

    let selected_searchables = create_selector(cx, move || {
        let mut selected = vec![];
        for searchable in searchables_ref.iter() {
            if *searchable.is_selected.get() {
                selected.push(searchable.clone());
            }
        }
        selected
    });

    searchables_ref.iter().for_each(|searchable| {
        create_effect(cx, move || {
            let is_selected = *searchable.is_selected.get();
            if is_selected {
                props
                    .selected_entities
                    .modify()
                    .insert(searchable.entity.clone());
            } else {
                props
                    .selected_entities
                    .modify()
                    .retain(|current| *current != searchable.entity);
            }
        });
    });

    let has_input_focus = create_signal(cx, false);
    let has_dialog_focus = create_signal(cx, false);

    let show_dialog = create_selector(cx, move || {
        *has_input_focus.get() || *has_dialog_focus.get()
    });

    let has_results = create_selector(cx, move || search_results.get().len() > 0);

    view! {cx,
        div(class="relative flex justify-center items-center text-gray-500") {
            Search(search_term=search_term, has_input_focus=has_input_focus, placeholder=props.placeholder)
            dialog(class="absolute z-20 top-full left-0 bg-white shadow-md rounded dark:bg-purple-700 dark:text-white",
                    open=*show_dialog.get(),
                    on:focusin= move |_| has_dialog_focus.set(true),
                    on:blur= move |_| has_dialog_focus.set(false),
                ) {
                (if *has_results.get() {
                    view! {cx,
                        ul(class="py-2 px-4 max-h-56 overflow-y-auto overflow-x-hidden divide-y") {
                            Indexed(
                                iterable=search_results,
                                view=move |cx, row| {
                                    let highlicht = create_selector(cx, move || {
                                        let name = row.entity.to_string();
                                        let search = search_term.get().to_ascii_lowercase();
                                        if let Some(index) = name.to_ascii_lowercase().find(&search) {
                                            let size = index + search.len();

                                            let prefix = name[..index].to_owned();
                                            let middle = name[index..size].to_owned();
                                            let suffix = name[size..].to_owned();
                                            (prefix, middle, suffix)
                                        } else {
                                            (name, "".to_string(), "".to_string())
                                        }
                                    });

                                    view! {cx,
                                        li(class="flex flex-col rounded hover:bg-gray-100 hover:dark:bg-purple-600 cursor-pointer",
                                            on:click=move |_| row.is_selected.set(!*row.is_selected.get())) {
                                            div(class="flex items-center justify-between my-2") {
                                                div(class="flex items-center") {
                                                    (highlicht.get().0)
                                                    span(class="font-bold") {
                                                        (highlicht.get().1)
                                                    }
                                                    (highlicht.get().2)
                                                }
                                                (if *row.is_selected.get() {
                                                    view! {cx,
                                                        span(class="w-6 h-6 bg-primary icon-[icon-park-solid--check-one]")
                                                    }
                                                } else {
                                                    view! {cx,
                                                        span(class="w-6 h-6 rounded-full border border-gray-300")
                                                    }
                                                })
                                            }
                                            hr(class="h-0.5 border-t-0 bg-neutral-100 opacity-100 dark:opacity-50 last:bg-transparent last:border-0")
                                        }
                                    }
                                },
                            )
                        }
                    }
                } else {
                    view! {cx,
                        p(class="text-center") {
                            "No results"
                        }
                    }
                })
            }
        }
        div(class="flex flex-wrap justify-center items-center mt-4") {
            Indexed(
                iterable= selected_searchables,
                view=move |cx, searchable| {
                    view!{cx,
                        div(class="flex items-center justify-center m-1 px-4 py-1 border-2 border-primary text-primary rounded-full") {
                            (searchable.entity.to_string())
                            span(class="w-4 h-4 ml-2 z-10 bg-primary icon-[material-symbols--close] cursor-pointer",
                                on:click=move |_| searchable.is_selected.set(false)
                            )
                        }
                    }
                }
            )
        }
        div(class="flex-grow")
    }
}
