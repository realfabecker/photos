import { Outlet } from "react-router";
import Header from "@pages/Layout/Header.tsx";
import Footer from "@pages/Layout/Footer.tsx";
import { useInjection } from "inversify-react";
import { IAuthService, Types } from "@core/ports/ports.ts";
import { Navigate } from "react-router-dom";
import { RoutesEnum } from "@core/domain/domain.ts";

export const PubLayout = () => {
  const service = useInjection<IAuthService>(Types.AuthService);
  if (service.isLoggedIn()) {
    return <Navigate to={RoutesEnum.Photos} />;
  }
  return (
    <div id="app">
      <Header />
      <main>
        <Outlet />
      </main>
      <Footer />
    </div>
  );
};

export const PrivLayout = () => {
  const service = useInjection<IAuthService>(Types.AuthService);
  if (!service.isLoggedIn()) {
    return <Navigate to={RoutesEnum.Login} />;
  }
  return (
    <div id="app">
      <Header />
      <main>
        <div className="container">
          <Outlet />
        </div>
      </main>
      <Footer />
    </div>
  );
};
