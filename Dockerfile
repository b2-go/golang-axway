FROM library/debian
ADD axwaymws /
CMD ["/axwaymws", "s.reagere.com:28000", "axway", "axwaymws"]
