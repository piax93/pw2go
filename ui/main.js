/**
 * Callback for click on service list item.
 * @param {string}   service - Service name.
 * @param {DOMEvent} event   - Click event.
 */
function onServiceClick(service, event) {
    for(var i = 0; i < event.target.classList.length; i++) {
        if(event.target.classList[i] == 'delete') return;
    }
    modal('Enter master password', function(res) {
        if(res.button === 'OK') {
            manager.get(service, res.master);
        }
    }, ['CANCEL', 'OK'], [{
        type: 'password',
        name: 'master',
        hint: 'Master Password'
    }]);
}

/**
 * Callback for click on service delete button.
 * @param {string}   service - Service name.
 * @param {DOMEvent} event   - Click event.
 */
function onDeleteClick(service, event) {
    modal('Do you really want to delete ' + service + '?', function(res) {
        if(res.button === 'YES') {
            manager.delete(service);
            modal('Service <i><b>' + service + '</b></i> deleted from database');
        }
    }, ['YES', 'NO']);
}

/**
 * Callback for click on ADD button
 * @param {DOMEvent} event - Click event.
 */
function onAddClick(event) {
    modal('Add new password', function(res) {
        if(res.button === 'ADD') {
            manager.add(res.service, res.password, res.master);
        }
    }, ['CANCEL', 'ADD'], [{
            type: 'text',
            name: 'service',
            hint: 'Service Name'
        }, {
            type: 'password',
            name: 'password',
            hint: 'Password'
        }, {
            type: 'password',
            name: 'master',
            hint: 'Master password'
        }
    ]);
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
        if(res.button === 'SET') {
            manager.setMaster(res.password);
            modal('First execution setup completed');
        } else {
            manager.die();
        }
    }, ['CANCEL', 'SET'], [{
            type: 'password',
            name: 'password',
            hint: 'Password',
        }]
    );
}

// Export functions
this.clearList = clearList;
this.setMaster = setMaster;
document.getElementById('addbtn').onclick = onAddClick;
