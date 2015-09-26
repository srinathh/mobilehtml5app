"use strict";

var TodoNew = React.createClass({
    displayName: "TodoNew",

    getInitialState: function getInitialState() {
        return { priority: 0 };
    },
    render: function render() {
        var thiscolor;
        if (this.state.priority == 0) {
            thiscolor = "grey";
        } else {
            thiscolor = "lightsalmon";
        }
        return React.createElement(
            "div",
            { className: "panel normalpanel" },
            React.createElement(
                "div",
                { className: "panel-body ", id: "createItemPanel" },
                React.createElement(
                    "table",
                    { width: "100%" },
                    React.createElement(
                        "tr",
                        null,
                        React.createElement(
                            "td",
                            { width: "80%" },
                            React.createElement("input", { type: "text", className: "form-control", id: "createItemText", placeholder: "Create a Todo" })
                        ),
                        React.createElement(
                            "td",
                            { className: "tdicon", width: "15%" },
                            React.createElement("span", { onClick: this.handleClickPriority, className: "glyphicon glyphicon-exclamation-sign todoicon", style: { "color": thiscolor } })
                        ),
                        React.createElement(
                            "td",
                            { className: "tdicon", width: "10%" },
                            React.createElement("span", { onClick: this.handleClickCreate, className: "glyphicon glyphicon-plus todoicon" })
                        )
                    )
                )
            )
        );
    },
    handleClickCreate: function handleClickCreate(evt) {
        var txt = $(React.findDOMNode(this)).find("#createItemText").val();
        if (txt == "") {
            return;
        }
        var impt = this.state.priority;
        createItem(txt, impt);
        $(React.findDOMNode(this)).find("#createItemText").val("");
        this.setState({ priority: 0 });
    },
    handleClickPriority: function handleClickPriority(evt) {
        if (this.state.priority == 0) {
            this.setState({ priority: 1 });
        } else {
            this.setState({ priority: 0 });
        }
    }
});

var Todo = React.createClass({
    displayName: "Todo",

    render: function render() {
        var colorclass;
        switch (this.props.item.Priority) {
            case 1:
                colorclass = "panel prioritypanel";
                break;
            default:
                colorclass = "panel normalpanel";
                break;
        }

        return React.createElement(
            "div",
            { className: colorclass },
            React.createElement(
                "div",
                { className: "panel-body " },
                React.createElement(
                    "table",
                    { width: "100%" },
                    React.createElement(
                        "tr",
                        null,
                        React.createElement(
                            "td",
                            null,
                            this.props.item.Text
                        ),
                        React.createElement(
                            "td",
                            { className: "tdicon", width: "10%" },
                            React.createElement("span", { onClick: this.handleClick, className: "glyphicon glyphicon-ok todoicon" })
                        )
                    )
                )
            )
        );
    },
    handleClick: function handleClick(evt) {
        deleteItem(this.props.item.ID);
    }
});

var TodoList = React.createClass({
    displayName: "TodoList",

    render: function render() {
        var self = this;
        var todos = this.props.items.map(function (item) {
            return React.createElement(Todo, { key: item.ID, item: item });
        });

        return React.createElement(
            "div",
            null,
            todos
        );
    },
    componentDidMount: function componentDidMount() {
        fetchItems();
    }
});

var deleteItem = function deleteItem(itemid) {
    $.ajax({
        url: "/items/" + itemid,
        type: "GET",
        success: function success() {
            fetchItems();
        },
        error: function error(xhr, status, errorThrown) {
            alert("Error deleting items - " + errorThrown);
        }

    });
};

var createItem = function createItem(txt, priority) {
    var item = { "ID": "newitem", "Text": txt, "Priority": priority };
    $.ajax({
        url: '/items/new',
        type: 'POST',
        data: { data: JSON.stringify(item) },
        success: function success(ret) {
            fetchItems();
        },
        error: function error(xhr, status, errorThrown) {
            alert("Error creating todo item");
        }
    });
};

var fetchItems = function fetchItems() {
    $.ajax({
        url: "/items",
        type: "GET",
        datatype: "json",
        success: function success(json) {
            var items = $.parseJSON(json);
            React.render(React.createElement(TodoNew, null), document.getElementById("todonew"));
            React.render(React.createElement(TodoList, { items: items }), document.getElementById("todolist"));
        },
        error: function error(xhr, status, errorThrown) {
            alert("Error fetching items - " + errorThrown);
        }

    });
};

$(function () {
    fetchItems();
});

