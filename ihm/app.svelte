<!--permet d'appeler le composant svelte dans le fichier html avec une simple balise-->
<svelte:options customElement="app-todo" />

<script context="module">
	export function customQueryEscape(str){
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

	export function answerResponse(text,statusCode) {
		try {
			if (statusCode == 403) {
				alert(`${text}reconnectez vous`);
				redirectTo("index.html");
			} else if (statusCode == 500) {
				alert(`${text}réessayez`);
			} else if (statusCode == 401) {
				alert(`${text}échec d'authentification`);
			}else {
				alert(`${text}`);
				console.log("statut d'erreur inattendu :", statusCode);
			}
		} catch (error) {
			console.error("erreur d'analyse de la réponse du serveur :", error);
		}
	}
	
	export async function logout(){
		try {
			const reponse = await fetch("/logout",{method: "GET"});
			if (!reponse.ok) {
				const errorData = await reponse.text();
				alert(errorData);
				throw new Error(`Erreur lors de la requête : ${reponse.status} ${reponse.statusText}`);
			}
			await reponse.text();
		} catch (error) {
			console.error(error.message);
		}
		redirectTo('index.html');
	}
</script>

<script>
	import { redirectTo } from './index.svelte';
	
	let todos = [];
	let nouvelleTache='';
	let selectedPriority = 2;

	function selectPriority(priority) {
		selectedPriority = priority;
	}

	function getPriorityColor(priority) {
		switch (priority) {
			case 3:
				return "rgba(255, 0, 0)";
			case 2:
				return "rgba(255, 255, 0)";
			case 1:
				return "rgba(0, 128, 0)";
			default:
				return "rgba(0, 0, 0)";
		}
	}

	function changePriority(item){
		if (item.priority >= 3){
			item.priority = 0;
			item.priority +=1;
		} else {
			item.priority += 1;
		}
		modify(item);
	}

	function add() {
		let todo = {
			text: nouvelleTache,
			priority: selectedPriority
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
				answerResponse(errorData,reponse.status);
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
			console.error(error.message);
		}
	}

	async function sendTodo(todo, route) {
		todo.text = customQueryEscape(todo.text)
		if(route=="add") {
			var url = `/${route}?text=${todo.text}&priority=${todo.priority}`;
		} else {
			url = `/${route}?id=${todo.id}&text=${todo.text}&priority=${todo.priority}`;
		}
		try {
			const reponse = await fetch(url,{method: "POST"});
			if (!reponse.ok) {
				const errorData = await reponse.text();
				answerResponse(errorData,reponse.status);
				throw new Error(`Erreur lors de la requête : ${reponse.status} ${reponse.statusText}`);
			}
			await reponse.json();
			getTodoList();
		} catch (error) {
			console.error(error.message);
		}
	}

	function handleKeydown(e, item) {
		if(e.key=="Enter") {
			if(item != null){
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
		on:keydown={e => handleKeydown(e, null)}
		/>
		<div class="priority">
			<button class="urgent {selectedPriority === 3 ? 'selectedurg' : ''}" on:click={() => selectPriority(3)}></button>
			<button class="prioritaire {selectedPriority === 2 ? 'selectedprio' : ''}" on:click={() => selectPriority(2)}></button>
			<button class="nonprioritaire {selectedPriority === 1 ? 'selectednonurg' : ''}" on:click={() => selectPriority(1)}></button>
		</div>
		<button class="ajout" disabled={nouvelleTache==""} on:click={add}>
			✔️
		</button>
	</div>
	
	<ul id="todo-list" class="todos">	
		{#each todos as item}
			<li>
				<input
					type="text"
					bind:value={item.text}
					on:keydown={e => handleKeydown(e, item)}
				/>
				<button
					class="priority-btn"
					on:click={() => changePriority(item)}
					style="background-color: {getPriorityColor(item.priority)}"
				>
				</button>
				<button class="app-btn" on:click={modify(item)}>
					✏️
				</button>
				<button class="app-btn" on:click={xclear(item)}>
					🗑️
				</button>
			</li>
		{/each}
	</ul>

</div>

<div class="bottom">

	<p>{remaining} tâches restantes !</p>
	<button class="disconnect" on:click={logout}>Se déconnecter</button>

</div>

<style>
	@import './style/style.css';
	.newTask {
		flex: 0.75;
		left: 12%;
		margin-bottom: 5%;
		margin-right: 1%;
		position: relative;
		font-size: large;
	}
	.ajout {
		height: 35px;
		width: 35px;
		border: none;
		border-radius: 5%;
		background-color: white;
		margin-left: 2%;
		position: relative;
	}
	.ajout:hover{
		background-color: rgba(146, 146, 146, 0.181);
		cursor: pointer;
	}
	.ajout:disabled{
		background-color: white;
		color: rgba(128, 128, 128, 0.836);
		cursor: default;
	}
	.todos input {
		font-size: large;
	}
	.priority-btn{
		position: relative;
		width: 15px;
		height: 12.5px;
		border: 2px solid rgba(0, 0, 0, 0.75);
		border-radius: 50%;
		margin: 5% 3%;
	}
	.priority-btn:hover{
		cursor: pointer;
	}
	.priority {
		display: inline;
		margin-left: 13%;
		justify-content: center;
	}
	.priority button{
		position: relative;
		width: 15px;
		height: 12.5px;
		border: 2px solid rgba(0, 0, 0, 0.75);
		border-radius: 50%;
		margin-right: 1%;
		margin-left: 1%;
		bottom: 10%;
	}
	.priority button:hover{
		cursor: pointer;
	}
	.urgent {
    	background-color: rgba(255, 0, 0, 0.25);
	}

	.prioritaire {
		background-color: rgba(255, 255, 0, 0.25);
	}
	.nonprioritaire {
		background-color: rgba(0, 128, 0, 0.25);
	}
	.selectedurg {
		border: 2px solid;
		border-color: black;
		background-color: red;
		transform: scale(1.2);
	}
	.selectedprio {
		border: 2px solid;
		border-color: black;
		background-color: yellow;
		transform: scale(1.2);
	}
	.selectednonurg {
		border: 2px solid;
		border-color: black;
		background-color: green;
		transform: scale(1.2);
	}
	.bottom{
    margin: auto;
    display: flex;
    justify-content: center;
    align-items: center;
	}
	.bottom p{
		font-size: large;
		margin: 1%;
		bottom: 7%;
	}
	.disconnect{
		width: 200px;
		height: 20px;
		margin: 1%;
		font-size: small;
		bottom: 1%;
		position: fixed;
		border: none;
		transition: transform 0.1s ease;
	}
	.disconnect:hover{
		background-color: rgba(146, 146, 146, 0.181);
		cursor: pointer;
	}
	.disconnect:active {
		transform: scale(0.95);
	}
	.app-btn {
		height: 35px;
		width: 35px;
		margin-left: 3%;
		transition: transform 0.2s ease;
	}
</style>
