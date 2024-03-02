/*
Package objectid implements ASN.1 Object Identifier types and methods.

# Features

  - Unsigned 128-bit numberForm support (allows for expressing registrations found below {joint-iso-itu-t(2) uuid(25)}, per ITU-T Rec. X.667.
  - Flexible index support, allowing interrogation through negative indices without the risk of panic
  - Convenient Leaf, Parent and Root index alias methods, wherever applicable
  - Ge, Gt, Le, Lt, Equal comparison methods for interacting with NumberForm instances

# License

The go-objectid package is available under the terms of the MIT license.

For further details, see the LICENSE file within the root of the source repository.

# NumberForm Maximum

Valid NumberForm instances may fall between the minimum decimal value of zero (0) and the maximum decimal value of 340,282,366,920,938,463,463,374,607,431,768,211,455 (three hundred forty undecillion and change). This ensures no panics occur when parsing valid UUID-based object identifiers, such as those found beneath joint-iso-itu-t(2) uuid(25) per [X.667](https://www.itu.int/rec/T-REC-X.667).

# Special Credit

A special thanks to Luke Champine for his excellent Uint128 package (found at https://github.com/lukechampine/uint128), which is incorporated within this package.
*/
package objectid
