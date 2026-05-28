type result = Ok(JSON.t) | Error(Err.t)

module Get = {
  let request = async (url): result => {
    let defaultErr = "Error processing API request"

    try {
      let resp = await Fetch.fetch(url)
      let json = await Response.json(resp)

      resp.ok ? Ok(json) : json->Err.parse(~defaultErr)->Error
    } catch {
    | JsExn(exn) => exn->Err.catch(~defaultErr)->Error
    }
  }
}
