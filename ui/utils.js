/**
 * Create DOM element.
 * @param {string} tag     - HTML tag name.
 * @param {string} id      - Element ID.
 * @param {array}  classes - List of classes for the element.
 * @return {DOMElement} Generated element.
 */
function createNode(tag, id, classes, inner) {
    classes = classes && 'length' in classes ? classes : [];
    var res = document.createElement(tag);
    if(id) res.id = id;
    for(var i = 0; i < classes.length; i++) {
        res.classList.add(classes[i])
    }
    if(inner) res.innerHTML = inner;
    return res;
}

/**
 * Create input element
 * @param {string} type  - Input type.
 * @param {string} name  - Input name.
 * @param {string} hint  - Placeholder.
 * @param {string} value - Input value.
 * @return {DOMElement} Input DOM element.
 */
function createInput(type, name, hint, value) {
    var classes = type === 'button' ? ['button', 'is-primary'] : [];
    var res = createNode('input', name, classes);
    res.type = type;
    res.name = name;
    if(hint) res.placeholder = hint;
    if(value) res.value = value;
    return res;
}

/**
 * Show popup
 * @param {string}   text    - Text to display in the popup.
 * @param {function} result  - Callback to call when button are pressed. The argument is an object with the button pressed and text fields values.
 * @param {array}    buttons - List of buttons to display.
 * @param {array}    fields  - List of text fields to display ({name, hint, type}).
 */
function modal(text, result, buttons, fields) {
    result = typeof result !== 'undefined' ? result : null;
    buttons = buttons && 'length' in buttons ? buttons : ['OK'];
    fields = fields && 'length' in fields ? fields : [];
    var body = createNode('div', null, ['modal', 'is-active']);
    body.appendChild(createNode('div', null, ['modal-background']));
    var content = createNode('div', null, ['modal-content', 'box']);
    body.appendChild(content);
    content.appendChild(createNode('p', null, null, text));
    for(var i = 0; i < fields.length; i++) {
        var f = fields[i];
        fields[i] = createInput(f.type, f.name, f.hint);
        content.appendChild(fields[i]);
    }
    var bcontainer = createNode('div', null, ['buttons', 'is-centered']);
    for(var i = 0; i < buttons.length; i++) {
        var b = buttons[i];
        var be = createInput('button', b, null, b);
        be.onclick = (function(bname, event) {
            if(typeof result === 'function') {
                var obj = {button: bname};
                for(var j = 0; j < fields[j]; j++) {
                    obj[fields[j].name] = fields[j].value;
                }
                result(obj);
            }
            document.getElementsByTagName('body')[0].removeChild(body);
        }).bind(null, b);
        bcontainer.appendChild(be);
    }
    content.appendChild(bcontainer);
    document.getElementsByTagName('body')[0].appendChild(body);
}

// Export functions
this.modal = modal;
this.createNode = createNode;
this.createInput = createInput;
