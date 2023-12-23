import "reflect-metadata";
import { Provider } from "react-redux";
import ReactDOM from "react-dom/client";
import { Provider as Container } from "inversify-react";
import App from "@pages/App.tsx";
import { store } from "@store/store.ts";
import { container } from "@core/container";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <Container container={container}>
    <Provider store={store}>
      <App />
    </Provider>
  </Container>
);
