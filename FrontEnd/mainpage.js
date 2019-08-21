function myCreateEvent() {
    var eventdata = document.getElementById("feedDetail").value

    var url = "http://localhost:8999/Events";
    var query = "{\"Shraddha\":\"" + eventdata + "\"}";
    const otherParams = {
        method: 'POST',
        headers: {
            "Content-Type": 'application/json; charset=utf-8',
        },
        body: query
    }
    fetch(url, otherParams)
        .then(data => {
           console.log(data)
        })
        .then(res => {
            console.log(res);
        })
        .then(error => {
            console.log(error);
        })
}