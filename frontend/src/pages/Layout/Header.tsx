const Nav = () => {
  return (
    <nav>
      <img src="/images/logo.svg" alt="Logotipo de câmera" className="logo" />

      <ul>
        <li className="active">Galeria</li>
      </ul>
    </nav>
  );
};

const Search = () => {
  return (
    <form className="search">
      <div className="input-wrapper">
        <label htmlFor="search">Pesquise por imagens e coleções</label>
        <input
          type="text"
          name="search"
          id="search"
          placeholder="Pesquise por imagens e coleções"
        />
      </div>
    </form>
  );
};

const Avatar = () => {
  return (
    <img
      src="https://gravatar.com/userimage/130158802/19f322df3acb8cd150b6c2cee26a6191.jpeg?size=256"
      alt="avatar image"
      className="avatar"
    />
  );
};

export default function Header() {
  return (
    <header className="container">
      <Nav />
      <Search />
      <Avatar />
    </header>
  );
}
