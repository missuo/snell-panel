import { useEffect, useState } from "react";
import { Button } from "@heroui/react";

export function ThemeToggle() {
  const [dark, setDark] = useState(
    () => document.documentElement.classList.contains("dark"),
  );

  useEffect(() => {
    const el = document.documentElement;
    el.classList.toggle("dark", dark);
    el.setAttribute("data-theme", dark ? "dark" : "light");
  }, [dark]);

  return (
    <Button variant="ghost" onPress={() => setDark((d) => !d)}>
      {dark ? "Light" : "Dark"}
    </Button>
  );
}
