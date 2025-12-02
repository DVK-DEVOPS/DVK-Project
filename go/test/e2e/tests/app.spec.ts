import { expect, test } from "@playwright/test";

test("homepage loads", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("textbox", { name: "Search..." })).toBeVisible();
});

test("search works", async ({ page }) => {
  await page.goto("/");
  await page.getByRole("textbox", { name: "Search..." }).fill("java");
  await page.getByRole("button", { name: "Search" }).click();
  await expect(
    page.getByText("Java Streams Examples v2 progurl2.dev-v2")
  ).toBeVisible();
});

test("navigation works", async ({ page }) => {
  await page.goto("/");
  await page.getByRole("link", { name: "Register" }).click();
  await expect(page).toHaveURL(/register/);
  await page.getByRole("link", { name: "Already have an account?" }).click();
  await expect(page).toHaveURL(/login/);
});

test("register validation", async ({ page }) => {
  await page.goto("/register");
  await page.getByRole("button", { name: "Sign up" }).click();
  await expect(page.getByRole("textbox", { name: "Search..." })).toBeVisible();
});

test("login flow", async ({ page }) => {
  await page.goto("/login");
  await page.locator('input[name="username"]').fill("test");
  await page.locator('input[name="password"]').fill("test");
  await page.getByRole("button", { name: "Log in" }).click();
  await expect(page.getByRole("textbox", { name: "Search..." })).toBeVisible();
});
