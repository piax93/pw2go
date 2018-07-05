/**
 * Callback for click on service list item.
 * @param {string}   service - Service name.
 * @param {DOMEvent} event   - Click event.
 */
function onServiceClick(service, event) {
    for(var i = 0; i < event.target.classList.length; i++) {
        if(event.target.classList[i] == 'delete') return;
    }
    modal(service);
}

/**
 * Callback for clickc on service delete button.
 * @param {string}   service - Service name.
 * @param {DOMEvent} event   - Click event.
 */
function onDeleteClick(service, event) {
    modal('Do you really want to delete ' + service + '?', function(res) {
        alert(res.button);
        if(res.button === 'YES') manager.delete(service);
    }, ['YES', 'NO']);
}

/**
 * Add service to list.
 * @param {string} service - Service name.
 */
this.addService = function addService(service) {
    var list = document.getElementById('servicelist');
    var li = createNode('li', null, ['service', 'box']);
    var cont = createNode('div', null, ['box'], service);
    var a = createNode('a', null, ['delete']);
    li.onclick = onServiceClick.bind(null, service);
    a.onclick = onDeleteClick.bind(null, service);
    li.appendChild(cont);
    li.appendChild(a);
    list.appendChild(li);
}

/**
 * Clear service list.
 */
function clearList() {
    var list = document.getElementById('servicelist')
    while (list.firstChild) {
        list.removeChild(list.firstChild);
    }
}

/**
 * Set master password.
 */
function setMaster() {
    modal('Set the master password', function(res){
        manager.setMaster(res.password);
    }, ['SET', 'CANCEL'], [{
            type: 'password',
            name: 'password',
            hint: 'Password',
        }]
    );
}

// Export functions
this.clearList = clearList;
this.setMaster = setMaster;
