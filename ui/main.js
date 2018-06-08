function addService(service) {
    var list = document.getElementById('servicelist');
    var li = document.createElement('li');
    var cont = document.createElement('div')
    var a = document.createElement('a');
    li.classList.add('service')
    li.classList.add('box');
    a.classList.add('delete');
    cont.classList.add('box');
    li.appendChild(cont);
    li.appendChild(a);
    list.appendChild(li);
    cont.innerHTML = service;
}
