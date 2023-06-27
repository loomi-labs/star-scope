use sycamore::prelude::*;

use crate::config::style;

#[component(inline_props)]
fn Button<'a, G: Html, F, S>(
    cx: Scope<'a>,
    children: Children<'a, G>,
    on_click: F,
    class: S,
) -> View<G>
where
    F: Fn() + 'a,
    S: Into<String>,
{
    let child_views = children.call(cx);
    let class_str = class.into();
    view!(
        cx,
        button(class=class_str, on:click=move |_| on_click()) {
            (child_views)
        }
    )
}

const BUTTON_BASE_CLASS: &str = "font-bold py-2 px-4 rounded";

pub enum ColorScheme<'a> {
    PrimaryGradient,
    Subtle,
    Custom(&'a str),
}

#[derive(Prop)]
pub struct ButtonProps<'a, G: Html, F> {
    children: Children<'a, G>,
    on_click: F,
    #[builder(default = ColorScheme::PrimaryGradient)]
    color: ColorScheme<'a>,
}

pub fn SolidButton<'a, G: Html, F>(cx: Scope<'a>, props: ButtonProps<'a, G, F>) -> View<G>
where
    F: Fn() + 'a,
{
    let color_class = match props.color {
        ColorScheme::PrimaryGradient => style::PRIMARY_TRANSIENT_BG_COLOR,
        ColorScheme::Subtle => "bg-purple-700 hover:bg-purple-800",
        ColorScheme::Custom(class) => class,
    };
    view!(
        cx,
        Button(
            children = props.children,
            on_click = props.on_click,
            class = format!("{} {}", BUTTON_BASE_CLASS, color_class)
        )
    )
}

pub fn OutlineButton<'a, G: Html, F>(cx: Scope<'a>, props: ButtonProps<'a, G, F>) -> View<G>
where
    F: Fn() + 'a,
{
    view!(
        cx,
        Button(children=props.children, on_click=props.on_click, class=format!("border border-primary-gradient-from transition-all hover:bg-primary-gradient-from {}", BUTTON_BASE_CLASS))
    )
}
