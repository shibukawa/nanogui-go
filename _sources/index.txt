.. NanoGUI.go documentation master file, created by
   sphinx-quickstart on Sat Dec  5 23:10:08 2015.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

NanoGUI.go
======================================

.. image:: https://godoc.org/github.com/shibukawa/nanogui.go?status.svg
   :target: https://godoc.org/github.com/shibukawa/nanogui.go

.. raw:: html

   <iframe src="/nanogui.go/demo/" width="900" height="540"></iframe>

`Full Screen </nanogui.go/demo/>`_

NanoGUI.go is apure golang implementation of `NanoGUI <https://github.com/wjakob/nanogui>`_. NanoGUI is a minimalistic GUI library for OpenGL.

NanoGUI.go can run on desktops and browsers (I didn't test on other environemnt).

It uses the following components:

* `OpenGL/WebGL library <https://github.com/goxjs/gl>`_
* `Cross platform glfw wrapper <https://github.com/goxjs/glfw>`_
* `NanoVGo <https://shibukawa.github.io/nanovgo/>`_

API Reference
---------------

See `GoDoc <https://godoc.org/github.com/shibukawa/nanogui.go>`_

Porting Status
------------------

.. list-table::
   :widths: 15 15 10 10
   :header-rows: 1

   - * Category
     * Classes
     * Finished
     * Status
   - * Non-Visual types
     * Screen
     * X
     * except cursor feature
   - * 
     * Widget
     * X
     *
   - *
     * GroupLayout, BoxLayout, GridLayout
     * X
     *
   - *
     * AdvancedGridLayout
     *
     *
   - *
     * Theme
     * X
     * Entype icons, Roboto fonts included
   - * Widgets
     * Window, Popup
     * X
     *
   - *
     * Label
     * X
     *
   - *
     * Button, ToolButton
     * X
     *
   - * 
     * PopupButton, ComboBox, ImagePanel, ColorPicker
     * X
     *
   - *
     * CheckBox
     * X
     *
   - *
     * Slider
     * X
     *
   - *
     * ProgressBar
     * X
     *
   - * 
     * TextBox, IntBox, FloatBox
     * X
     * It supports Emacs key bind like MacOS X, IME not supported.
   - *
     * ImageView, Graph
     * X
     *
   - *
     * VScrollPanel
     * X
     *
   - *
     * ColorWheel
     * X
     *
   - *
     * MessageDialog
     *
     *
   - * Utitility
     * FormHelper
     *
     *
   - * OpenGL Helper
     * GLShader, GLFramebuffer, Arcball
     *
     *

Author
---------------

* `Yoshiki Shibukawa <https://github.com/shibukawa>`_

License
----------

NanoGUI.go is released under zlib license. But the following contents are under other licenses:

* Roboto fonts: Apatch 2 license
* Entypo icon: Creative Commons 4.0 BY-SA

Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`

