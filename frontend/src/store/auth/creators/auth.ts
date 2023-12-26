import { createAsyncThunk } from "@reduxjs/toolkit";
import { Container } from "inversify";
import { IAuthService, Types } from "@core/ports/ports";
import { NavigateFunction, Location } from "react-router";
import { RoutesEnum } from "@core/domain/domain";

export const getActionAuthLogin = createAsyncThunk(
  "auth/login",
  async (
    {
      email,
      password,
      navigate,
      location,
    }: {
      email: string;
      password: string;
      navigate: NavigateFunction;
      location: Location;
    },
    { extra }
  ) => {
    const container = (<any>extra).container as Container;
    const authService = container.get<IAuthService>(Types.AuthService);
    await authService.login({ email, password });
    navigate(RoutesEnum.Photos + location.search);
  }
);

export const getActionAuthLogout = createAsyncThunk(
  "auth/logout",
  async ({ navigate }: { navigate: NavigateFunction }, { extra }) => {
    const container = (<any>extra).container as Container;
    const authService = container.get<IAuthService>(Types.AuthService);
    authService.logout();
    navigate(RoutesEnum.Login);
  }
);
