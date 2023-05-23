use log::debug;
use wasm_bindgen::JsValue;

pub fn add_or_update_query_params(key: &str, value: &str) {
    let window = web_sys::window().unwrap();
    let history = window.history().unwrap();
    let document = window.document().unwrap();
    let url = document.url().unwrap();
    let segments = url.split("?").collect::<Vec<_>>();
    let stripped_url = segments[0];
    let mut query_params = querystring::querify("");
    if segments.len() > 1 {
        query_params = querystring::querify(segments[1]);
    }
    let index = query_params.iter().position(|(k, _)| k.to_string() == key);
    debug!("query_params: {:?}", query_params);
    if let Some(index) = index {
        query_params.remove(index);
    }
    if !value.is_empty() {
        query_params.push((key, value));
    }
    let new_query_params = querystring::stringify(vec![(key, value)]).trim_end_matches("&").to_string();
    let new_url_string = format!("{}?{}", stripped_url, new_query_params);
    history.replace_state_with_url(&JsValue::null(), "", Some(&new_url_string)).unwrap();
}

pub fn get_query_param(key: &str) -> Option<String> {
    let window = web_sys::window().unwrap();
    let document = window.document().unwrap();
    let url = document.url().unwrap();
    let query_params = querystring::querify(url.as_str());
    let index = query_params.iter().position(|(k, _)| k.to_string() == key);
    if let Some(index) = index {
        Some(query_params[index].1.to_string())
    } else {
        None
    }
}