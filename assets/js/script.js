document.addEventListener("DOMContentLoaded", e => {
    parseDate()
})

function parseDate() {
    let lis = document.getElementsByTagName("li")
    let start = "Start: "
    let end = "End: "
    for (let i = 0; i < lis.length; i++) {
        let d = lis[i].innerHTML
        if (d.includes(start)) {
            newDate = new Date(d.slice(7, 36))
            if (isValidDate(newDate)) {
                lis[i].innerHTML = start + newDate.toUTCString()
            }
        }
        if (d.includes(end)) {
            newDate = new Date(d.slice(5, 34))
            if (isValidDate(newDate)) {
                lis[i].innerHTML = end + newDate.toUTCString()
            }
        }
    }
}

function isValidDate(d) {
    return d instanceof Date && !isNaN(d);
}