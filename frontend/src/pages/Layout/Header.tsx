import { useInjection } from "inversify-react";
import { IAuthService, Types } from "@core/ports/ports.ts";

const Nav = (opts: { isLoggedIn: boolean }) => {
  return (
    <nav>
      <img src="/images/logo.svg" alt="Logotipo de câmera" className="logo" />

      {opts.isLoggedIn && (
        <ul>
          <li className="active">Galeria</li>
        </ul>
      )}
    </nav>
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
  const service = useInjection<IAuthService>(Types.AuthService);
  const isLoggedIn = service.isLoggedIn();
  return (
    <header className={`container ${isLoggedIn ? "private" : "public"}`}>
      <Nav isLoggedIn={isLoggedIn} />
      {service.isLoggedIn() && <Avatar />}
    </header>
  );
}
