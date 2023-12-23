import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import { configureStore } from "@reduxjs/toolkit";
import { ActionStatus } from "@core/domain/domain";
import { container } from "@core/container";
import photoSlice from "@store/photos/reducers/photos";

export interface State<T = any> {
  data?: T;
  status: ActionStatus;
  error?: { message: string };
}

export const store = configureStore({
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      thunk: { extraArgument: { container } },
    }),
  reducer: {
    photos: photoSlice,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
