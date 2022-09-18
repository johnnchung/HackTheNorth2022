import "./App.css";
import { TextArea } from "./components";
import { CookieMonster, DonaldTrump, BenShapiro, Penguin, Giraffe, Logo } from "./components/images"

function App() {
  return (
    <div className="App">
      <CookieMonster/>
      <DonaldTrump/>
      <BenShapiro/>
      <Penguin/>
      <Giraffe/>
      <Logo/>
      <TextArea></TextArea>
    </div>
  );
}

export default App;
