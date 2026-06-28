import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { StatusBadge } from "./status-badge";

describe("StatusBadge", () => {
  it("todo は『未着手』と表示する", () => {
    render(<StatusBadge status="todo" />);
    expect(screen.getByText("未着手")).toBeInTheDocument();
  });

  it("doing は『進行中』と表示する", () => {
    render(<StatusBadge status="doing" />);
    expect(screen.getByText("進行中")).toBeInTheDocument();
  });

  it("done は『完了』と表示する", () => {
    render(<StatusBadge status="done" />);
    expect(screen.getByText("完了")).toBeInTheDocument();
  });
});
