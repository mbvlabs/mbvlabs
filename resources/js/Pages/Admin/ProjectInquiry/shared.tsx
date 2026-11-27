import { useState } from 'react'
import { Button, buttonVariants } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from '@/Components/ui/alert-dialog'

export const button = buttonVariants({ variant: 'outline' })

export function Confirm({ label, title, description, confirmation, disabled, actionLabel, onConfirm }: { label: string; title: string; description: string; confirmation?: string; disabled?: boolean; actionLabel: string; onConfirm: () => void }) {
  const [value, setValue] = useState('')
  return <AlertDialog onOpenChange={open => { if (!open) setValue('') }}>
    <AlertDialogTrigger asChild><Button type="button" variant="outline" disabled={disabled}>{label}</Button></AlertDialogTrigger>
    <AlertDialogContent>
      <AlertDialogHeader><AlertDialogTitle>{title}</AlertDialogTitle><AlertDialogDescription>{description}</AlertDialogDescription></AlertDialogHeader>
      {confirmation && <div><Label htmlFor={`confirm-${label}`}>Type {confirmation} to confirm</Label><Input id={`confirm-${label}`} className="mt-2" autoComplete="off" value={value} onChange={e => setValue(e.target.value)} /></div>}
      <AlertDialogFooter><AlertDialogCancel>Cancel</AlertDialogCancel><AlertDialogAction disabled={disabled || !!confirmation && value !== confirmation} onClick={onConfirm}>{actionLabel}</AlertDialogAction></AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
}
