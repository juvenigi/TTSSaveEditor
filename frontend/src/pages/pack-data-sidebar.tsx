import {Button} from "@/components/ui/button.tsx";
import {
  Drawer,
  DrawerClose,
  DrawerContent, DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger
} from "@/components/ui/drawer.tsx";

export default function PackDataSidebar() {

  return (
    <Drawer>
      <DrawerTrigger>Open PackData</DrawerTrigger>
      <DrawerContent>
        <DrawerHeader>
          <DrawerTitle>Pack Data</DrawerTitle>
          <DrawerDescription>ooh yeah</DrawerDescription>
        </DrawerHeader>
        <DrawerFooter>
          <Button>Submit</Button>
          <DrawerClose>
            <Button variant="outline">Cancel</Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  )
}