function onServiceClick(service, event) {
    for(var i = 0; i < event.target.classList.length; i++)
        if(event.target.classList[i] == 'delete') return;
    alert(service);
}

function onDeleteClick(service, event) {
    if(confirm("Do you really want to delete " + service + "?"))
        manager.delete(service);
}

function modal(text, button, result) {

}

this.addService = function addService(service) {
    var list = document.getElementById('servicelist');
    var li = document.createElement('li');
    var cont = document.createElement('div');
    var a = document.createElement('a');
    li.classList.add('service');
    li.classList.add('box');
    li.onclick = onServiceClick.bind(null, service);
    a.classList.add('delete');
    a.onclick = onDeleteClick.bind(null, service);
    cont.classList.add('box');
    li.appendChild(cont);
    li.appendChild(a);
    list.appendChild(li);
    cont.innerHTML = service;
}

this.clearList = function clearList() {
    var list = document.getElementById('servicelist')
    while (list.firstChild) {
        list.removeChild(list.firstChild);
    }
}

this.setMaster = function setMaster() {
    modal('Set the master password', 'SET', function(res){
        manager.setMaster(res);
    });
}