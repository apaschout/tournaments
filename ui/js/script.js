document.addEventListener("DOMContentLoaded", e => {
    limitDropdown(4)
    parseDate()
})

document.getElementById("form-changename").addEventListener("submit", e => {
    e.preventDefault()
    console.log("got here")
    $.ajax({
        url: "http://127.0.0.1:5000/",
        type: "post",
        data: $('#form-changename').serialize(),
        success: () => {
            window.location.reload()
        }
    })
})

function limitDropdown(limit) {
    let trns = document.getElementsByClassName("form-dropdown-tournaments")
    let plrs = document.getElementsByClassName("form-dropdown-players")
    for (let i = limit; i < trns.length; i++) {
        trns[i].style.display = "none"
    }
    for (let i = limit; i < plrs.length; i++)[
        plrs[i].style.display = "none"
    ]
}

function parseDate() {
    let lis = document.getElementsByTagName("li")
    let start = "Start: "
    let end = "End: "
    for (let i = 0; i < lis.length; i++) {
        let d = lis[i].innerHTML
        if (d.includes(start)) {
            console.log(d.slice(7, 36))
            lis[i].innerHTML = start + new Date(d.slice(7, 36)).toUTCString()
        }
        if (d.includes(end)) {
            console.log(d.slice(5, 34))
            lis[i].innerHTML = end + new Date(d.slice(5, 34)).toUTCString()
        }
    }
}

function subForm() {
    console.log($('#form-changename'))
    $('#form-changename').submit(e => {
        e.preventDefault()
        console.log("got here")
        $.ajax({
            url: "http://127.0.0.1:5000/",
            type: "post",
            data: $('#form-changename').serialize(),
            success: () => {
                window.location.reload()
            }
        })
    })
}