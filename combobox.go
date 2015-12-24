package nanogui

type ComboBox struct {
	PopupButton
	callback      func(int)
	items         []string
	shortItems    []string
	selectedIndex int
}

func NewComboBox(parent Widget, items ...[]string) *ComboBox {
	var itemsParam []string
	var shortItemsParam []string
	switch len(items) {
	case 0:
	case 1:
		itemsParam = items[0]
	case 2:
		itemsParam = items[0]
		shortItemsParam = items[1]
	default:
		panic("NewComboBox can accept upto 2 extra parameters (items, shortItems)")
	}
	var index int
	if len(items) == 0 {
		index = -1
	}
	combobox := &ComboBox{
		selectedIndex: index,
	}
	// init PopupButton member
	combobox.chevronIcon = IconRightOpen
	combobox.SetIconPosition(ButtonIconLeftCentered)
	combobox.SetFlags(ToggleButtonType | PopupButtonType)
	parentWindow := parent.FindWindow()
	combobox.popup = NewPopup(parentWindow.Parent(), parentWindow)
	combobox.popup.SetSize(320, 250)
	InitWidget(combobox, parent)
	combobox.SetItems(itemsParam, shortItemsParam)
	return combobox
}

func (c *ComboBox) SelectedIndex() int {
	return c.selectedIndex
}

func (c *ComboBox) SetSelectedIndex(i int) {
	if len(c.shortItems) == 0 {
		return
	}
	children := c.PopupButton.Popup().Children()
	if c.selectedIndex > -1 {
		children[c.selectedIndex].(*Button).SetPushed(false)
	}
	if i < 0 || i >= len(c.items) {
		c.selectedIndex = -1
		c.SetCaption("")
	} else {
		children[i].(*Button).SetPushed(true)
		c.selectedIndex = i
		c.SetCaption(c.shortItems[i])
	}
}

func generateCallback(c *ComboBox, popup *Popup, i int) func() {
	return func() {
		c.selectedIndex = i
		c.SetCaption(c.shortItems[i])
		c.SetPushed(false)
		popup.SetVisible(false)
		if c.callback != nil {
			c.callback(i)
		}
	}
}

func (c *ComboBox) SetItems(items []string, shortItems ...[]string) {
	var shortItemsParam []string
	switch len(shortItems) {
	case 0:
	case 1:
		shortItemsParam = shortItems[0]
	default:
		panic("ComboBox.SetItems can accept only one extra parameter (shortItems)")
	}
	if len(shortItemsParam) == 0 {
		shortItemsParam = items
	}
	if len(items) != len(shortItemsParam) {
		panic("ComboBox.SetItems can accept only same length string lists as items and shortItems.")
	}
	c.items = items
	c.shortItems = shortItemsParam
	if c.selectedIndex < 0 || c.selectedIndex >= len(c.items) {
		c.selectedIndex = -1
	}
	popup := c.Popup()
	for popup.ChildCount() > 0 {
		popup.RemoveChildByIndex(popup.ChildCount() - 1)
	}
	popup.SetLayout(NewGroupLayout(10))
	for i, item := range items {
		button := NewButton(popup, item)
		button.SetFlags(RadioButtonType)
		button.SetCallback(generateCallback(c, popup, i))
	}
	c.SetSelectedIndex(c.selectedIndex)
}

func (c *ComboBox) Items() []string {
	return c.items
}

func (c *ComboBox) ShortItems() []string {
	return c.shortItems
}

func (c *ComboBox) SetCallback(callback func(int)) {
	c.callback = callback
}

func (c *ComboBox) String() string {
	return c.StringHelper("ComboBox", c.caption)
}
