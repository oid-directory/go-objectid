/*
Package objectid implements ASN.1 Object Identifier types and methods.

# Features

• Unsigned 128-bit arc support (i.e.: such as the registrations found below {joint-iso-itu-t(2) uuid(25)})

• Flexible index support, allowing interrogation through negative indices without the risk of panic

• Convenient Leaf, Parent and Root index alias methods, wherever allowed

# Uint128 Support

Unsigned 128-bit integer support for individual NumberForm values is made possible due to the private incorporation of Luke Champine's awesome Uint128 type, which manifests here through instances of the package-provided NumberForm type.

Valid NumberForm instances may fall between the minimum decimal value of zero (0) and the maximum decimal value of 340,282,366,920,938,463,463,374,607,431,768,211,455 (three hundred forty undecillion and change). This ensures no panics occur when parsing valid UUID-based object identifiers.
*/
package objectid
