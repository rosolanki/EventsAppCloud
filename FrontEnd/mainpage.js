function myCreateEvent() {
    var userinfo = document.getElementById("userDetail").value
    var eventdata = document.getElementById("feedDetail").value

    var url = "http://localhost:8999/Events";
    var query = "{\"user\":\"" + userinfo + "\", \"feed\":\"" + eventdata + "\"}";
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
            window.alert(data.status)
        })
        .then(res => {
            console.log(res);
        })
        .then(error => {
            console.log(error);
        })
}

function getEvent() {
    $("#feedDisplay").empty()
    var userID = document.getElementById("userGetID").value

    var url = "http://localhost:8999/Events?user=" + userID;
    const otherParams = {
        method: 'GET',
        headers: {
            "Content-Type": 'application/json; charset=utf-8',
        }
    }
    fetch(url, otherParams)
        .then(response => response.json())
        .then(data => {
            console.log(data)
            $("#feedDisplay").append("<h5>"+ data.feed +"</h5>")
        })
        .then(error => {
            console.log(error);
        })
}