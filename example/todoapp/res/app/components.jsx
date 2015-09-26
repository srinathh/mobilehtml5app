        var TodoNew = React.createClass({
            getInitialState: function() {
                return {priority: 0};
            },
            render:function(){
                var thiscolor;
                if(this.state.priority == 0){
                    thiscolor = "grey"
                }else{
                    thiscolor = "lightsalmon"
                }
                return(
                    <div className="panel normalpanel">
            	    	<div className="panel-body " id="createItemPanel">
                    		<table width="100%">
                    			<tr>
                                    <td width="80%"><input type="text" className="form-control" id="createItemText" placeholder="Create a Todo" /></td>
        				            <td className="tdicon" width="15%"><span onClick={this.handleClickPriority} className="glyphicon glyphicon-exclamation-sign todoicon" style={{"color":thiscolor}}></span></td>
        				            <td className="tdicon" width="10%"><span onClick={this.handleClickCreate} className="glyphicon glyphicon-plus todoicon"></span></td>
            	        		</tr>
                    		</table>
            	    	</div>
                    </div>
                );
            },
            handleClickCreate:function(evt){
                var txt = $(React.findDOMNode(this)).find("#createItemText").val();
                if(txt == ""){
                    return;
                }
                var impt = this.state.priority;
                createItem(txt, impt);
                $(React.findDOMNode(this)).find("#createItemText").val("");
                this.setState({priority:0});
            },
            handleClickPriority:function(evt){
                if(this.state.priority ==0){
                    this.setState({priority:1});
                } else{
                    this.setState({priority:0});
                }
            }
        });

        var Todo = React.createClass({
        	render:function(){
        		var colorclass;
        		switch(this.props.item.Priority){
        		case 1:
        			colorclass = "panel prioritypanel"
        			break;
        		default:
        			colorclass = "panel normalpanel"
        			break;
        		}

        		return(
        			<div className={colorclass}>
        	    	    <div className="panel-body ">
                				<table width="100%">
                					<tr>
        	        					<td>{this.props.item.Text}</td>
        	        					<td className="tdicon" width="10%"><span onClick={this.handleClick} className="glyphicon glyphicon-ok todoicon"></span></td>
        	        				</tr>
                				</table>
        	    	    </div>
        		    </div>
        		);
        	},
        	handleClick:function(evt){
                deleteItem(this.props.item.ID);
            }
        });


        var TodoList = React.createClass({
        	render:function(){
        		var self = this;
        		var todos = this.props.items.map(function(item){
        			return(
        				<Todo key = {item.ID} item={item}/>
        			)
        		});

        		return(
        			<div>
        				{todos}
        			</div>
        		);
        	},
            componentDidMount:function(){
                fetchItems();
            }
        });

        var deleteItem = function(itemid){
        	$.ajax({
        		url:"/items/"+itemid,
        		type:"GET",
        		success:function(){
                    fetchItems();
                },
        		error:function(xhr, status, errorThrown){
        			alert("Error deleting items - "+ errorThrown)
        		}

        	});
        };

        var createItem=function(txt, priority){
            var item = {"ID":"newitem","Text":txt,"Priority":priority}
        	$.ajax({
        		url: '/items/new',
        		type: 'POST',
        		data: {data: JSON.stringify(item)},
        		success:function(ret){
        			fetchItems();
        		},
        		error:function(xhr, status, errorThrown){
        			alert("Error creating todo item")
        		}
        	});
        };

        var fetchItems = function(){
        	$.ajax({
        		url:"/items",
        		type:"GET",
        		datatype:"json",
        		success:function(json){
        			var items = $.parseJSON(json);
                    React.render(
                        <TodoNew />,
                        document.getElementById("todonew")
                    );
        			React.render(
        				<TodoList items={items} />,
        				document.getElementById("todolist")
        			);
        		},
        		error:function(xhr, status, errorThrown){
        			alert("Error fetching items - "+ errorThrown)
        		}

        	});
        };

        $(function(){
            fetchItems();
        });
