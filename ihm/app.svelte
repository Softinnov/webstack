<!--permet d'appeler le composant svelte dans le fichier html avec une simple balise-->
<svelte:options customElement="app-todo" />

<script>
	let todos = [];
	let nouvelleTache='';

	function customQueryEscape(str){
		const uriStr = encodeURIComponent(str);
		const finalQueryStr = uriStr
		.replace(/!/g, '%21')
		.replace(/'/g, '%27')
		.replace(/\(/g, '%28')
        .replace(/\)/g, '%29')
		.replace(/\*/g, '%2A')
		.replace(/%20/g,'+');
    	return finalQueryStr;
	}

	function add() {
		let todo = {
			text: nouvelleTache
		}
		try{
			sendTodo(todo,"add");
		} catch (error) {
			console.error(`Erreur lors de la connection au serveur : ${error.message}`);
		}
		nouvelleTache='';
	}

	function xclear(item) {
		try {
			sendTodo(item,"delete");
		} catch (error) {
			console.error(`Erreur lors de la connection au serveur : ${error.message}`);
		}
	}

	function modify(item) {
		try {
			sendTodo(item,"modify");
		} catch (error) {
			console.error(`Erreur lors de la connection au serveur : ${error.message}`);
		}
	}

	async function getTodoList() {
		const url = `/todos`	
		try {
			const reponse = await fetch(url,{method: "GET"});
			if (!reponse.ok) {
				const errorData = await reponse.text();
				alert(errorData);
				throw new Error(`Erreur lors de la requête : ${reponse.status} ${reponse.statusText}`);
			}
			const result = await reponse.json();
			if(result != null){
				todos = result;
			}
			if(result == null){
				todos = [];
			}
		} catch (error) {
			console.error(`Une erreur s'est produite : ${error.message}`);
		}
	}

	async function sendTodo(todo, route) {
		todo.text = customQueryEscape(todo.text)
		if(route=="add") {
			var url = `/${route}?text=${todo.text}`;
		} else {
			url = `/${route}?id=${todo.id}&text=${todo.text}`;
		}
		try {
			const reponse = await fetch(url,{method: "POST"});
			if (!reponse.ok) {
				const errorData = await reponse.text();
				alert(errorData);
				throw new Error(`Erreur lors de la requête : ${reponse.status} ${reponse.statusText}`);
			}
			const result = await reponse.json();
			getTodoList();
		} catch (error) {
			console.error(`Une erreur s'est produite : ${error.message}`);
		}
	}

	function handleKeydown(e, item, action) {
		if(e.key=="Enter") {
			// console.log("Enter");
			if(action=="modify"){
				console.log("on est la");
				modify(item);
			}
			else {
				add();
			}
		}
	}

	$: remaining = todos.length;	
	
	getTodoList();

</script>

<div class="centered">
	
	<h1>My TodoList</h1>

	<div>
		<input  
		class="newTask"
		type="text" 
		placeholder="Quoi d'autre?"
		bind:value={nouvelleTache}
		on:keydown={handleKeydown}
		/>
		<button class="ajout" disabled={nouvelleTache==""} on:click={add}>
			✔️
		</button>
	</div>
	
	<ul id="todo-list" class="todos">	
		{#each todos as item}
			<li>
				<button class="button" on:click={modify(item)}>
					✏️
				</button>
				<button class="button" on:click={xclear(item)}>
					🗑️
				</button>

				<input
					id="todo"
					type="text"
					bind:value={item.text}
					on:keydown={handleKeydown(item, "modify")}
				/>
			</li>
		{/each}
	</ul>

	<p>{remaining} tâches restantes !</p>

</div>

<style>
	h1{
		font-size: 70px;
	}
	p{
		font-size: large;
		left: 45%;
		bottom: 5%;
		position: fixed;
	}
	.button {
		height: 35px;
		width: 35px;
		border: none;
		border-radius: 5%;
		background-color: white;
		margin-right: 3%;
	}
	.button:hover{
		background-color: rgba(146, 146, 146, 0.381);
		cursor: pointer;
	}
	.centered {
		width: 25em;
		margin: auto;
		display:grid;
	}
	.newTask {
		left: 12%;
		margin-bottom: 15%;
		margin-right: 1%;
		position: relative;
	}
	.todos {
		margin-top: 10%;
	}
	.ajout {
		height: 35px;
		width: 35px;
		border: none;
		border-radius: 5%;
		background-color: white;
		margin-left: 13%;
		position: relative;
	}
	.ajout:hover{
		background-color: rgba(146, 146, 146, 0.381);
		cursor: pointer;
	}
	.ajout:disabled{
		background-color: white;
		color: rgba(128, 128, 128, 0.836);
		cursor: default;
	}
	ul {
		max-height: 15em;
		overflow: hidden;
		overflow-y: visible;
	}
	li {
		display: flex;
		margin: 3%;
	}
	input[type="text"] {
		flex: 1;
		padding: 0.5em;
		margin: -0.2em 0;
		border: none;
		font-size: large;
	}
</style>
