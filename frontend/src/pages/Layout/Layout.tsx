import { ReactNode } from "react";
import Header from "@pages/Layout/Header.tsx";
import Footer from "@pages/Layout/Footer.tsx";

export default function Layout(opts: { children: ReactNode }) {
  return (
    <div id="app">
      <Header />
      <main>
        <div className="container">{opts.children}</div>
      </main>
      <Footer />
    </div>
  );
}
