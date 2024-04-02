<svelte:options customElement="login-todo" />

<script>
    let email = '';
    let password = '';
  
    const handleSubmit = () => {
        console.log("Email:", email);
        console.log("Password:", password);
    };
</script>

  
<div>
    <h2>Connexion</h2>
    <form on:submit|preventDefault={handleSubmit}>
        <label>
        Email:
            <input type="email" bind:value={email} required>
        </label>
  
        <label>
        Mot de passe:
            <input type="password" bind:value={password} required>
        </label>
  
        <button type="submit">Se connecter</button>
    </form>
</div>
  