"use client";

import Link from "next/link";
import { useActionState } from "react";

import type { AuthState } from "@/app/(auth)/actions";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

type Props = {
  title: string;
  submitLabel: string;
  action: (prev: AuthState, formData: FormData) => Promise<AuthState>;
  altHref: string;
  altLabel: string;
};

export function AuthForm({
  title,
  submitLabel,
  action,
  altHref,
  altLabel,
}: Props) {
  const [state, formAction, pending] = useActionState<AuthState, FormData>(
    action,
    {},
  );

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <form action={formAction}>
        <CardContent className="space-y-4">
          {state.error && (
            <p role="alert" className="text-sm text-destructive">
              {state.error}
            </p>
          )}
          <div className="space-y-2">
            <Label htmlFor="email">メールアドレス</Label>
            <Input
              id="email"
              name="email"
              type="email"
              required
              autoComplete="email"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">パスワード</Label>
            <Input
              id="password"
              name="password"
              type="password"
              required
              minLength={8}
              autoComplete="current-password"
            />
          </div>
        </CardContent>
        <CardFooter className="mt-4 flex-col items-stretch gap-3">
          <Button type="submit" disabled={pending}>
            {pending ? "送信中..." : submitLabel}
          </Button>
          <p className="text-center text-sm text-muted-foreground">
            <Link href={altHref} className="underline">
              {altLabel}
            </Link>
          </p>
        </CardFooter>
      </form>
    </Card>
  );
}
